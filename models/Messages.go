package models

import "time"

type Messages struct {
	ID        int       `json:"id"`
	UserId    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
