package main

import (
	"fmt"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/cfg"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/handler"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	authHandler := handler.NewAuthorizeHandler()
	authHandler.RegisterRoutes(r)

	loginHandler := handler.NewLoginHandler()
	loginHandler.RegisterRoutes(r)

	signupHandler := handler.NewSignupHandler()
	signupHandler.RegisterRoutes(r)

	fmt.Println("Server on...")
	http.ListenAndServe(cfg.Env().HttpServerPort, r)

}
