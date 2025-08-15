# Ronin-C2

**Educational / Lab-Only Command & Control Framework**  
Built in Go with a Vue 3 dashboard for studying C2 tasking, execution, and monitoring patterns in authorized environments.

⚠ **Disclaimer**  
This project is intended **solely** for authorized security research, lab simulations, and training.  
Running it against systems without explicit permission is **illegal** and **unethical**.  
The authors accept no liability for misuse.

---

## ✨ Features
- Go-based server with HTTP API and WebSocket dashboard
- Per-agent task queues
- Agent polling with configurable sleep/jitter
- API key authentication (via `RONIN_C2_API_KEY`)
- Web UI for real-time monitoring
- Educational, lab-safe defaults

---

## 🗂 Project Structure
```
cmd/server/      → C2 server entrypoint  
cmd/agent/       → Agent entrypoint  
internal/        → Shared logic (router, middleware, ws, agent)  
ui/              → Vue 3 dashboard (Vite)  
```

---

## 🔧 Setup

### Prerequisites
- Go 1.20+
- Node.js 18+ (for UI)
- Git, Docker (optional)

### 1. Clone & configure
```bash
git clone git@github.com:zerotrace0x/ronin-c2.git
cd ronin-c2
export RONIN_C2_API_KEY="supersecretkey"
```

### 2. Run server
```bash
cd cmd/server
go run main.go
```
Server: `127.0.0.1:7777` (API)

### 3. Run agent
```bash
cd cmd/agent
export RONIN_C2_AGENT_ID="agent-123"
go run main.go
```

### 4. Queue a command
```bash
curl -sX POST http://127.0.0.1:7777/command   -H "Authorization: Bearer $RONIN_C2_API_KEY"   -H "Content-Type: application/json"   -d '{"agent_id":"agent-123","command":"whoami"}' | jq
```

### 5. Check results
```bash
curl -s "http://127.0.0.1:7777/results?agent_id=agent-123"   -H "Authorization: Bearer $RONIN_C2_API_KEY" | jq
```

## 📜 License
MIT — see [LICENSE](LICENSE)

**Author:** [zerotrace0x](https://github.com/zerotrace0x)
