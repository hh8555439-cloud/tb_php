package controller

import (
	"net/http"
	"time"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// 清除session_id Cookie
	expiration := time.Now().Add(-24 * time.Hour)

	// 清除jwt_token Cookie
	jwtCookie := &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: expiration,
		Path:    "/",
	}
	http.SetCookie(w, jwtCookie)

	// 重定向到登录页
	http.Redirect(w, r, "/login", http.StatusFound)
}
