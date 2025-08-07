package controller

import (
	"comment_demo/database"
	"comment_demo/models"
	"encoding/json"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
	"log"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("your-secret-key"))

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	User    struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Role     string `json:"role"`
	} `json:"user"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	response := LoginResponse{
		Success: false,
		Message: "",
	}

	session, err := store.Get(r, "session-name")
	if err != nil {
		response.Message = "会话初始化失败"
		json.NewEncoder(w).Encode(response)
		return
	}

	if r.Method == "POST" {
		var request LoginRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			response.Message = "请求数据解析失败"
			json.NewEncoder(w).Encode(response)
			return
		}

		// 验证输入
		if request.Username == "" || request.Password == "" {
			response.Message = "用户名和密码不能为空"
			json.NewEncoder(w).Encode(response)
			return
		}

		// 查询数据库
		var user models.User
		result := database.DB.Where("username = ? AND password = ?", request.Username, request.Password).First(&user)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				response.Message = "用户名或密码错误"
			} else {
				log.Printf("数据库查询错误: %v", result.Error)
				response.Message = "系统错误，请稍后再试"
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		// 设置会话变量
		session.Values["user_id"] = user.ID
		session.Values["username"] = request.Username
		session.Values["role"] = string(user.Role)
		err = session.Save(r, w)
		if err != nil {
			response.Message = "会话保存失败"
			json.NewEncoder(w).Encode(response)
			return
		}

		// 返回成功响应
		response.Success = true
		response.Message = "登录成功"
		response.User.ID = user.ID
		response.User.Username = request.Username
		response.User.Role = string(user.Role)
		json.NewEncoder(w).Encode(response)
	}
}
