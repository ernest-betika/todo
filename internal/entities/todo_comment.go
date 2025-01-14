package entities

import "time"

type TodoComment struct {
	Identifier
	Comment   string    `json:"comment"`
	TodoID    int64     `json:"todo_id"`
	CreatedAt time.Time `json:"created_at"`
}
