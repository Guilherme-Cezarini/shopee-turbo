package main

import (
	"context"
	"fmt"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/cfg"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/controllers"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/dependencies"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/repository/access_token"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/repository/user"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/service/utils"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog/log"

	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	deps, _ := dependencies.LoadGlobalDependencies(context.Background())
	utilsService := utils.NewUtilsService()
	CSFRKey, err := utilsService.GenerateRandomKey(32)
	CSRF := csrf.Protect(CSFRKey)

	sessionStore := sessions.NewCookieStore(CSFRKey)
	sessionStore.Options.HttpOnly = true
	sessionStore.Options.Secure = cfg.Env().SessionSecure
	sessionStore.Options.MaxAge = cfg.Env().SessionMaxAge * 3600

	userRepo := user.NewUserRepository(deps.Mongo.Database)
	accessTokenRespo := access_token.NewAccessTokenRepository(deps.Mongo.Database)

	r := mux.NewRouter()

	loginHandler := controllers.NewLoginHandler(userRepo, deps.HashKit, sessionStore)
	loginHandler.RegisterRoutes(r)

	signupHandler := controllers.NewSignupHandler(userRepo, deps.HashKit)
	signupHandler.RegisterRoutes(r)

	dashboardHandler := controllers.NewDashboardHandler(sessionStore, userRepo, accessTokenRespo)
	dashboardHandler.RegisterRoutes(r)

	fmt.Println("Server on...")
	err = http.ListenAndServe(cfg.Env().HttpServerPort, CSRF(r))
	if err != nil {
		log.Fatal().Err(err).Msg("Error starting server")
		return
	}

}
