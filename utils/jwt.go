package utils

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
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

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "未授权，缺少 Authorization 头", http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "未授权，Authorization 格式错误", http.StatusUnauthorized)
			return
		}
		tokenStr := parts[1]

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "未授权，token无效", http.StatusUnauthorized)
			return
		}

		// 将用户信息放入上下文，以便后续处理
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
