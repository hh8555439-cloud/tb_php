package route

import (
	g "comment_demo/controller"
	"comment_demo/database"
	"comment_demo/service"
	"net/http"
)

func SetupRoutes() {
	http.HandleFunc("/login", g.LoginHandler)
	http.HandleFunc("/register", g.RegisterHandler)

	commentService := service.NewCommentService(database.DB) // 你自己实现
	commentController := g.NewCommentController(commentService)

	http.HandleFunc("/api", commentController.ApiHandler)

}
