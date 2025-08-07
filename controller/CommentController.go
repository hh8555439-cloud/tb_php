package controller

import (
	"comment_demo/models"
	"comment_demo/service"
	"encoding/json"
	"net/http"
	"strconv"
)

type CommentController struct {
	service *service.CommentService
}

func NewCommentController(svc *service.CommentService) *CommentController {
	return &CommentController{service: svc}
}

// ApiHandler 统一处理 /api 路由
func (cc *CommentController) ApiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	action := r.URL.Query().Get("action")

	switch action {
	case "get_comments":
		cc.getComments(w, r)
	case "add_comment":
		cc.addComment(w, r)
	case "get_user":
		cc.getUser(w, r)
	// 你可以继续添加其他 action
	default:
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    1,
			"message": "Invalid action",
		})
	}
}

func (cc *CommentController) getUser(w http.ResponseWriter, r *http.Request) {
	// 假设你用 session 包管理登录状态
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	session, err := store.Get(r, "session-name")
	if err != nil {
		// session 获取失败，认为未登录
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 1,
			"data": nil,
		})
		return
	}

	userID, ok1 := session.Values["user_id"].(uint)
	username, ok2 := session.Values["username"].(string)
	role, ok3 := session.Values["role"].(string)

	if !ok1 || !ok2 || !ok3 {
		// session 中没有登录信息
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 1,
			"data": nil,
		})
		return
	}

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

func (cc *CommentController) getComments(w http.ResponseWriter, r *http.Request) {
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

func (cc *CommentController) addComment(w http.ResponseWriter, r *http.Request) {
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
