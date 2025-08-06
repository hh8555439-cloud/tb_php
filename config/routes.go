package route

import (
	g "comment_demo/go"
	"net/http"
)

func SetupRoutes() {
	http.HandleFunc("/login", g.LoginHandler)
	http.HandleFunc("/register", g.RegisterHandler)
}
