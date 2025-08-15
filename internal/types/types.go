package types

import "time"

type Command struct {
	ID        string    `json:"id"`
	AgentID   string    `json:"agent_id"`
	Command   string    `json:"command"`
	QueuedAt  time.Time `json:"queued_at"`
}

type Result struct {
	ID        string    `json:"id"`
	AgentID   string    `json:"agent_id"`
	CommandID string    `json:"command_id"`
	Stdout    string    `json:"stdout"`
	Stderr    string    `json:"stderr"`
	Code      int       `json:"code"`
	EndedAt   time.Time `json:"ended_at"`
}
