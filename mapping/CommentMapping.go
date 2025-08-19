package mapping

import (
	"comment_demo/models"
	"errors"
	"gorm.io/gorm"
	"time"
)

type CommentMapping struct {
	db *gorm.DB
}

type MessageDTO struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func NewCommentMapping(db *gorm.DB) *CommentMapping {
	return &CommentMapping{db: db}
}

// UserInfo 结构体定义
type UserInfo struct {
	ID       int    `gorm:"column:id"`
	Username string `gorm:"column:username"`
}

func (m *CommentMapping) GetMessages() ([]MessageDTO, error) {
	var messages []MessageDTO
	// 执行查询（使用预加载获取关联用户）
	stmt := m.db.Table("messages").
		Select("messages.id", "messages.user_id", "users.username", "messages.content", "messages.created_at").
		Joins("LEFT JOIN users ON messages.user_id = users.id").
		Order("messages.created_at DESC")
	result := stmt.Scan(&messages)
	if result.Error != nil {
		return nil, result.Error
	}
	return messages, nil
}

// AddComment 添加评论
func (m *CommentMapping) AddComment(comment models.Comment) (int, error) {
	if err := m.db.Create(&comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

// AddMessage 添加留言
func (m *CommentMapping) AddMessage(message models.Messages) (int, error) {
	if err := m.db.Create(&message).Error; err != nil {
		return 0, err
	}
	return message.ID, nil
}

// GetCommentsByGoods 获取商品所有评论
func (m *CommentMapping) GetCommentsByGoods(goodsID int) ([]models.Comment, error) {
	var comments []models.Comment
	err := m.db.Where("goods_id = ?", goodsID).Order("create_time ASC").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// GetUserInfo 获取用户信息（模拟，实际应从用户服务获取）
func (m *CommentMapping) GetUserInfo(userID int) (UserInfo, error) {
	var user UserInfo
	err := m.db.Table("users").Select("id, username").Where("id = ?", userID).Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return UserInfo{
				ID:       userID,
				Username: "未知用户",
			}, nil
		}
		return UserInfo{}, err
	}
	return user, nil
}

func (m *CommentMapping) DeleteComment(commentId int) error {
	result := m.db.Table("comments").Where("id = ?", commentId).Delete(nil)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *CommentMapping) DeleteMessage(messageId int) error {
	result := m.db.Table("messages").Where("id = ?", messageId).Delete(nil)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
