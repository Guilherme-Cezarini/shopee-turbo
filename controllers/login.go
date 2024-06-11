package controllers

import (
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/middlewares"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/models/requests"
	user2 "github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/repository/user"
	authorization_ml "github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/service/authorization-ml"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/service/hash"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/service/session"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog/log"
	"net/http"
	"text/template"
)

type LoginHandler struct {
	validator              *validator.Validate
	hashKit                hash.HashingInterface
	userRepo               user2.UserRepository
	sessionService         session.InternalSessionService
	sessionStore           *sessions.CookieStore
	authorizationMlService authorization_ml.AuthorizationService
}

var loginFront = template.Must(template.ParseGlob("front-end/views/login/*"))

var errorLoginOrPasswordInvalid = map[string]any{"Error": "Email ou senha inválidos."}

func NewLoginHandler(userRepo user2.UserRepository, hashKit hash.HashingInterface, sessionStore *sessions.CookieStore) *LoginHandler {
	return &LoginHandler{
		validator:              validator.New(),
		hashKit:                hashKit,
		userRepo:               userRepo,
		sessionService:         session.NewSessionService(),
		sessionStore:           sessionStore,
		authorizationMlService: authorization_ml.NewAuthorizationService(),
	}
}

func (handler *LoginHandler) executeFront(w http.ResponseWriter, page string, data any) error {
	err := loginFront.ExecuteTemplate(w, page, data)
	if err != nil {
		return err
	}

	return nil
}

func (handler *LoginHandler) RegisterRoutes(r *mux.Router) {

	r.Handle("/login", middlewares.AlreadyLoggedInMiddleware(handler.sessionStore)(http.HandlerFunc(handler.Index))).Methods("GET")
	r.HandleFunc("/sign-in", handler.Signin).Methods("POST")
}

// generateViewData generate the data to view
func (handler *LoginHandler) generateViewData(dataError any, login requests.Login, r *http.Request) map[string]interface{} {
	data := map[string]interface{}{
		"Error":          dataError,
		"Login":          login,
		csrf.TemplateTag: csrf.TemplateField(r),
	}
	return data

}

func (handler *LoginHandler) Index(w http.ResponseWriter, r *http.Request) {
	dataView := map[string]any{
		csrf.TemplateTag: csrf.TemplateField(r),
	}
	err := loginFront.ExecuteTemplate(w, "index", dataView)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (handler *LoginHandler) Signin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		login := requests.Login{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}
		validationError := handler.validator.Struct(login)
		if validationError != nil {
			dataToView := handler.generateViewData(validationError, login, r)
			err := handler.executeFront(w, "index", dataToView)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		passwordHashed := handler.hashKit.HashField(login.Password)
		user, err := handler.userRepo.GetUserByEmail(login.Email)
		if err != nil {
			dataToView := handler.generateViewData(errorLoginOrPasswordInvalid, login, r)
			err = handler.executeFront(w, "index", dataToView)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return

			}
			return
		}
		if passwordHashed != user.Password {
			dataToView := handler.generateViewData(errorLoginOrPasswordInvalid, login, r)
			err = handler.executeFront(w, "index", dataToView)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		err = handler.sessionService.CreateUserSession(handler.sessionStore, "session", user.Id, w, r)
		if err != nil {
			dataToView := handler.generateViewData(errorLoginOrPasswordInvalid, login, r)
			err = handler.executeFront(w, "index", dataToView)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		url := handler.authorizationMlService.GenerateUrlToMLAuthorization()
		log.Debug().Msgf("[DashBoardHandler] Redirect user %s to ML authorization with URL: %s", user.Id, url)
		http.Redirect(w, r, url, http.StatusSeeOther)
		return
	}
	err := loginFront.ExecuteTemplate(w, "index", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
