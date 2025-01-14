package entities

import (
	"time"
)

type Todo struct {
	Identifier
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	CompletedAt *time.Time `json:"completed_at"`
	Timestamps
}

type TodoList struct {
	Todos []*Todo `json:"todos"`
}
