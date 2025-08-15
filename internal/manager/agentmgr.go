package manager

import (
	"sync"

	"github.com/zerotrace0x/ronin-c2/internal/queue"
	"github.com/zerotrace0x/ronin-c2/internal/types"
)

type AgentManager struct {
	mu        sync.RWMutex
	commands  map[string]*queue.Queue[types.Command] // per-agent command queues
	results   map[string][]types.Result              // per-agent result history (latest first)
	maxResult int
}

func NewAgentManager(maxResult int) *AgentManager {
	return &AgentManager{
		commands:  make(map[string]*queue.Queue[types.Command]),
		results:   make(map[string][]types.Result),
		maxResult: maxResult,
	}
}

func (m *AgentManager) ensure(agentID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.commands[agentID]; !ok {
		m.commands[agentID] = queue.New[types.Command]()
	}
	if _, ok := m.results[agentID]; !ok {
		m.results[agentID] = make([]types.Result, 0, m.maxResult)
	}
}

func (m *AgentManager) Enqueue(agentID string, cmd types.Command) {
	m.ensure(agentID)
	m.commands[agentID].Enqueue(cmd)
}

func (m *AgentManager) Dequeue(agentID string) (types.Command, bool) {
	m.ensure(agentID)
	return m.commands[agentID].Dequeue()
}

func (m *AgentManager) AppendResult(agentID string, res types.Result) {
	m.ensure(agentID)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.results[agentID] = append([]types.Result{res}, m.results[agentID]...)
	if len(m.results[agentID]) > m.maxResult {
		m.results[agentID] = m.results[agentID][:m.maxResult]
	}
}

func (m *AgentManager) Results(agentID string, limit int) []types.Result {
	m.mu.RLock()
	defer m.mu.RUnlock()
	rs := m.results[agentID]
	if limit > 0 && len(rs) > limit {
		return rs[:limit]
	}
	return rs
}
