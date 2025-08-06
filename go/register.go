package _go

import (
	"comment_demo/database"
	"comment_demo/models"
	"encoding/json"
	"gorm.io/gorm"
	"log"
	"net/http"

	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Response struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

type RegisterRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var requestBody RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Message: "无效的请求格式",
			Errors:  map[string]string{"global": "请求体必须是有效的JSON"},
		})
		return
	}

	response := Response{
		Success: false,
		Message: "",
		Errors:  make(map[string]string),
	}

	// 输入验证
	if strings.TrimSpace(requestBody.Username) == "" {
		response.Errors["username"] = "用户名不能为空"
	}
	if strings.TrimSpace(requestBody.Email) == "" || !strings.Contains(requestBody.Email, "@") {
		response.Errors["email"] = "邮箱格式无效"
	}
	if requestBody.Password == "" {
		response.Errors["password"] = "密码不能为空"
	}
	if requestBody.Password != requestBody.ConfirmPassword {
		response.Errors["confirm_password"] = "两次密码不一致"
	}

	if len(response.Errors) == 0 {
		// 使用事务确保原子性
		err = database.DB.Transaction(func(tx *gorm.DB) error {
			// 检查用户名和邮箱唯一性
			var count int64
			if err := tx.Model(&models.User{}).
				Where("username = ? OR email = ?", requestBody.Username, requestBody.Email).
				Count(&count).Error; err != nil {
				log.Println("检查唯一性时出错:", err)
				return err
			}

			if count > 0 {
				response.Message = "用户名或邮箱已存在"
				response.Errors["global"] = "用户名或邮箱已存在"
				return nil
			}

			// 创建用户记录
			user := models.User{
				Username: requestBody.Username,
				Email:    requestBody.Email,
				Password: requestBody.Password, // 实际应用中应该使用哈希密码
			}

			if err := tx.Create(&user).Error; err != nil {
				log.Println("创建用户记录时出错:", err)
				return err
			}

			response.Success = true
			response.Message = "注册成功"
			return nil
		})

		if err != nil {
			log.Println("数据库操作失败:", err)
			response.Message = "系统错误，请稍后再试"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
