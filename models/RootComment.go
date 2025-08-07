package models

import "time"

type RootComment struct {
	ID         int             `json:"id"`
	Content    string          `json:"content"`
	User       User            `json:"user"`
	CreateTime time.Time       `json:"create_time"`
	Answers    []AnswerComment `json:"answers"`
	GoodsID    int             `json:"goods_id"`
}

type AnswerComment struct {
	ID         int       `json:"id"`
	Content    string    `json:"content"`
	User       User      `json:"user"`
	ToUser     *User     `json:"to_user,omitempty"`
	CreateTime time.Time `json:"create_time"`
	RootID     int       `json:"root_id"`
	ToAnswerID *int      `json:"to_answer_id,omitempty"`
	GoodsID    int       `json:"goods_id"`
}
