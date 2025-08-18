package models

import "time"

var messages []struct {
	ID   uint `json:"id"`
	User struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
