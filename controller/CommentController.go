package controller

import (
	"comment_demo/models"
	"comment_demo/service"
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strconv"
)

type CommentController struct {
	service *service.CommentService
}

func NewCommentController(svc *service.CommentService) *CommentController {
	return &CommentController{service: svc}
}

func (cc *CommentController) GetUser(w http.ResponseWriter, r *http.Request) {
	// 从Cookie中获取jwt_token
	cookie, err := r.Cookie("token")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 1,
			"data": nil,
		})
		return
	}

	// 解析JWT令牌
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil // 替换为你的JWT密钥
	})

	if err != nil || !token.Valid {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 1,
			"data": nil,
		})
		return
	}

	// 提取JWT中的用户数据
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 1,
			"data": nil,
		})
		return
	}

	userID := uint(claims["user_id"].(float64))
	username := claims["username"].(string)
	role := claims["role"].(string)

	// 返回用户信息
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"id":       userID,
			"username": username,
			"role":     role,
		},
	})
}

func (cc *CommentController) GetComments(w http.ResponseWriter, r *http.Request) {
	goodsIDStr := r.URL.Query().Get("goods_id")
	goodsID, err := strconv.Atoi(goodsIDStr)
	if err != nil || goodsID <= 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    1,
			"message": "invalid goods_id",
		})
		return
	}

	comments, err := cc.service.GetGoodsComments(goodsID)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    1,
			"message": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"data": comments,
	})
}

func (cc *CommentController) AddComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    1,
			"message": "Method not allowed",
		})
		return
	}

	var req struct {
		Content    string `json:"content"`
		UserID     int    `json:"user_id"`
		GoodsID    int    `json:"goods_id"`
		ToUserID   *int   `json:"to_user_id"`
		RootID     *int   `json:"root_id"`
		ToAnswerID *int   `json:"to_answer_id"`
		Type       string `json:"type"`
	}

	// 解析 JSON 请求体
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    1,
			"message": "Invalid request body",
		})
		return
	}

	comment := models.Comment{
		Content:    req.Content,
		UserID:     req.UserID,
		GoodsID:    req.GoodsID,
		ToUserID:   req.ToUserID,
		RootID:     req.RootID,
		ToAnswerID: req.ToAnswerID,
		Type:       req.Type,
	}

	id, err := cc.service.AddComment(comment)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    1,
			"message": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"data": id,
	})
}
