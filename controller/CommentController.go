package controller

import (
	"comment_demo/models"
	"comment_demo/service"
	"comment_demo/utils"
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strconv"
	"time"
)

type CommentController struct {
	service *service.CommentService
}

func NewCommentController(svc *service.CommentService) *CommentController {
	return &CommentController{service: svc}
}

func (cc *CommentController) GetMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := cc.service.GetMessages()
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    1,
			"message": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"data": messages,
	})
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
		CreateTime: time.Now(),
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

func (cc *CommentController) AddMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    1,
			"message": "Method not allowed",
		})
		return
	}

	var req struct {
		UserID  int    `json:"user_id"`
		Content string `json:"content"`
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

	message := models.Messages{
		UserId:    req.UserID,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	id, err := cc.service.AddMessage(message)
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

func (cc *CommentController) DeleteComment(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(utils.UserContextKey)
	if claims == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    1,
			"message": "无授权信息",
		})
		return
	}
	if claims.(jwt.MapClaims)["role"] != "admin" {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    1,
			"message": "无权限",
		})
		return
	}
	commentId, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    1,
			"message": "参数错误",
		})
		return
	}
	err = nil
	err = cc.service.DeleteComment(commentId)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    1,
			"message": err.Error(),
		})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"data": nil,
	})
}

func (cc *CommentController) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(utils.UserContextKey)
	if claims == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    1,
			"message": "无授权信息",
		})
		return
	}
	if claims.(jwt.MapClaims)["role"] != "admin" {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    1,
			"message": "无权限",
		})
		return
	}
	messageId, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    1,
			"message": "参数错误",
		})
		return
	}
	err = nil
	err = cc.service.DeleteMessage(messageId)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    1,
			"message": err.Error(),
		})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"data": nil,
	})
}
