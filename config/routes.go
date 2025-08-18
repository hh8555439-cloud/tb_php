package config

import (
	g "comment_demo/controller"
	"comment_demo/database"
	"comment_demo/service"
	"comment_demo/utils"
	"net/http"
)

func SetupRoutes() {
	http.HandleFunc("/login", g.LoginHandler)
	http.HandleFunc("/register", g.RegisterHandler)
	http.HandleFunc("/logout", g.LogoutHandler)

	commentService := service.NewCommentService(database.DB) // 你自己实现
	commentController := g.NewCommentController(commentService)

	http.HandleFunc("/api/get_comments", commentController.GetComments)
	http.HandleFunc("/api/get_user", commentController.GetUser)
	http.HandleFunc("/api/get_messages", commentController.GetMessages)

	// 需要保护
	http.Handle("/api/add_comment", utils.JWTAuthMiddleware(http.HandlerFunc(commentController.AddComment)))

}
