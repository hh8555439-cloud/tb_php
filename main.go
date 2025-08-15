package main

import (
	"comment_demo/config"
	"comment_demo/utils"
	"log"
	"net/http"
)

func main() {
	// 应用全局CORS中间件
	config.SetupRoutes()
	log.Println("服务器启动，监听端口 :8080")
	log.Fatal(http.ListenAndServe(":8080", utils.EnableCORS(http.DefaultServeMux)))
}
