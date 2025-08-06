package main

import (
	config "comment_demo/config"
	"log"
	"net/http"
)

func main() {
	config.SetupRoutes()
	log.Println("服务器启动，监听端口 :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
