package utils

import "net/http"

// 在主路由配置前添加全局CORS中间件
func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 允许前端来源
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:63342")

		// 允许认证凭证（如Cookie）
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// 允许的请求头
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 处理OPTIONS预检请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
