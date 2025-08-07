package mapping

import (
	"comment_demo/models"
	"errors"
	"gorm.io/gorm"
)

type CommentMapping struct {
	db *gorm.DB
}

func NewCommentMapping(db *gorm.DB) *CommentMapping {
	return &CommentMapping{db: db}
}

// UserInfo 结构体定义
type UserInfo struct {
	ID   int
	Name string
}

// AddComment 添加评论
func (m *CommentMapping) AddComment(comment models.Comment) (int, error) {
	if err := m.db.Create(&comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
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
				ID:   userID,
				Name: "未知用户",
			}, nil
		}
		return UserInfo{}, err
	}
	return user, nil
}
