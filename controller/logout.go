package controller

import (
	"net/http"
	"time"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// 清除会话变量（在Go中通常使用JWT或自定义会话管理）
	// 示例：假设使用cookie存储会话ID
	expiration := time.Now().Add(-24 * time.Hour)
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "",
		Expires: expiration,
		Path:    "/",
	}
	http.SetCookie(w, cookie)

	// 重定向到登录页
	http.Redirect(w, r, "/login", http.StatusFound)
}
