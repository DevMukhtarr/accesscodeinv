package routes

import (
	"net/http"

	"github.com/devmukhtarr/accesscodeinv/controllers"
	"github.com/devmukhtarr/accesscodeinv/middlewares"
)

func NewUser() {
	http.HandleFunc("/user/signin", controllers.SignIn)
	http.HandleFunc("/user/signup", controllers.SignUp)
	http.Handle("/token/new", middlewares.CheckToken(http.HandlerFunc(controllers.GetAccessToken)))
}
