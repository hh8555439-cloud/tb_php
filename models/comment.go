package models

import "time"

type Comment struct {
	ID         int       `json:"id"`
	Content    string    `json:"content"`
	UserID     int       `json:"user_id"`
	GoodsID    int       `json:"goods_id"`
	ToUserID   *int      `json:"to_user_id,omitempty"`
	RootID     *int      `json:"root_id,omitempty"`
	ToAnswerID *int      `json:"to_answer_id,omitempty"`
	Type       string    `json:"type"`
	CreateTime time.Time `json:"create_time"`
	// 数据库中能存null值的可以用指针形式
}
