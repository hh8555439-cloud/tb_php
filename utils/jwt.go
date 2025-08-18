package utils

import (
	"context"
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

var jwtKey = []byte("secret")

// Claims 定义 JWT 负载结构，根据需要定义字段
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// 生成 JWT Token
func GenerateJWT(userID uint, username, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // 例如 24 小时有效期
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "comment_demo",
			Subject:   "user token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

type key int

const UserContextKey key = 0

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Unauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(ErrorResponse{
		Code:    401,
		Message: message,
	})
}

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			Unauthorized(w, "未授权，缺少 token")
			return
		}
		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			Unauthorized(w, "未授权，无效的令牌")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			Unauthorized(w, "未授权，无法解析令牌内容")
			return
		}

		// 将用户信息放入上下文，以便后续处理
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
