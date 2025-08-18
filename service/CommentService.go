package service

import (
	mp "comment_demo/mapping"
	mo "comment_demo/models"
	"errors"
	"gorm.io/gorm"
	"sort"
	"time"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

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

// 转换为所需的响应格式
type ResponseMessage struct {
	ID   uint `json:"id"`
	User struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentService struct {
	model *mp.CommentMapping
}

func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{
		model: mp.NewCommentMapping(db),
	}
}

// GetGoodsComments 获取某商品的所有评论，包含根评论和回复
func (c *CommentService) GetGoodsComments(goodsID int) ([]RootComment, error) {
	// 从mapping层获取所有评论
	comments, err := c.model.GetCommentsByGoods(goodsID)
	if err != nil {
		return nil, err
	}

	// 缓存用户信息，避免重复查询
	userCache := make(map[int]User)
	// 根评论map，key是根评论ID
	rootComments := make(map[int]RootComment)

	for _, comment := range comments {
		// 缓存评论用户信息
		if _, ok := userCache[comment.UserID]; !ok {
			userInfo, err := c.model.GetUserInfo(comment.UserID)
			if err != nil {
				return nil, err
			}
			userCache[comment.UserID] = User{
				ID:       userInfo.ID,
				Username: userInfo.Username,
			}
		}

		if comment.Type == "root" {
			// 根评论
			rootComments[comment.ID] = RootComment{
				ID:         comment.ID,
				Content:    comment.Content,
				User:       userCache[comment.UserID],
				CreateTime: comment.CreateTime,
				Answers:    make([]AnswerComment, 0),
				GoodsID:    comment.GoodsID,
			}
		} else {
			// 回复评论，先确认根评论存在
			if comment.RootID == nil {
				// 逻辑上不应该出现，防御性编程
				continue
			}
			rootComment, ok := rootComments[*comment.RootID]
			if !ok {
				// 可能根评论不在当前列表，跳过或根据业务决定
				continue
			}

			// 缓存被回复用户信息
			var toUser *User
			if comment.ToUserID != nil {
				if _, ok := userCache[*comment.ToUserID]; !ok {
					userInfo, err := c.model.GetUserInfo(*comment.ToUserID)
					if err != nil {
						return nil, err
					}
					userCache[*comment.ToUserID] = User{
						ID:       userInfo.ID,
						Username: userInfo.Username,
					}
				}
				u := userCache[*comment.ToUserID]
				toUser = &u
			}

			answer := AnswerComment{
				ID:         comment.ID,
				Content:    comment.Content,
				User:       userCache[comment.UserID],
				ToUser:     toUser,
				CreateTime: comment.CreateTime,
				RootID:     *comment.RootID,
				ToAnswerID: comment.ToAnswerID,
				GoodsID:    comment.GoodsID,
			}

			// 添加回复到对应根评论的Answers切片
			rootComment.Answers = append(rootComment.Answers, answer)
			rootComments[*comment.RootID] = rootComment
		}
	}

	// 转成切片返回
	result := make([]RootComment, 0, len(rootComments))
	for _, rc := range rootComments {
		result = append(result, rc)
	}

	// 按根评论创建时间倒序排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreateTime.After(result[j].CreateTime)
	})

	return result, nil
}

// AddComment 添加评论，参数直接用models.Comment结构体
func (c *CommentService) AddComment(comment mo.Comment) (int, error) {
	// 如果是回复评论，必须指定RootID
	if comment.Type == "answer" && comment.RootID == nil {
		return 0, errors.New("回复评论必须指定root_id")
	}

	return c.model.AddComment(comment)
}

// AddMessage 添加留言，参数直接用models.Message结构体
func (c *CommentService) AddMessage(message mo.Messages) (int, error) {
	return c.model.AddMessage(message)
}

func (c *CommentService) GetMessages() ([]ResponseMessage, error) {
	messages, err := c.model.GetMessages()
	if err != nil {
		// 处理错误
		return nil, err
	}
	responseData := make([]ResponseMessage, len(messages))
	for i, msg := range messages {
		responseData[i] = ResponseMessage{
			ID: msg.ID,
			User: struct {
				ID   uint   `json:"id"`
				Name string `json:"name"`
			}{
				ID:   msg.UserID,
				Name: msg.Username,
			},
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
		}
	}
	return responseData, nil
}
