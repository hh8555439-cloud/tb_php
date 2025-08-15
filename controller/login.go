package controller

import (
	"comment_demo/database"
	"comment_demo/models"
	"comment_demo/utils"
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
	response := LoginResponse{
		Success: false,
		Message: "",
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

		// 生成JWT
		token, err := utils.GenerateJWT(user.ID, request.Username, string(user.Role))
		if err != nil {
			response.Message = "令牌生成失败"
			json.NewEncoder(w).Encode(response)
			return
		}

		// 返回成功响应
		response.Success = true
		response.Message = "登录成功"
		response.User.ID = user.ID
		response.User.Username = request.Username
		response.User.Role = string(user.Role)
		// 设置Cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			MaxAge:   86400, // 有效期1天
			HttpOnly: true,  // 防止XSS攻击
			Secure:   false, // 如果使用HTTPS，设为true
			SameSite: http.SameSiteLaxMode,
		})
		json.NewEncoder(w).Encode(response)
	}
}
