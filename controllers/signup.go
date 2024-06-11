package controllers

import (
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/models/user"
	user2 "github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/repository/user"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/service/hash"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"text/template"
)

type SignupHandler struct {
	validator *validator.Validate
	userRepo  user2.UserRepository
	hashKit   hash.HashingInterface
}

var signupFront = template.Must(template.ParseGlob("front-end/views/signup/*"))

func NewSignupHandler(userRepo user2.UserRepository, hashKit hash.HashingInterface) *SignupHandler {
	return &SignupHandler{
		validator: validator.New(),
		userRepo:  userRepo,
		hashKit:   hashKit,
	}
}

func (handler *SignupHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/sign-up", handler.Index).Methods("GET")
	r.HandleFunc("/sign-up/create", handler.Signup).Methods("POST")
}

func (handler *SignupHandler) Index(w http.ResponseWriter, r *http.Request) {
	err := signupFront.ExecuteTemplate(w, "index", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// executeFront call the execute front
func (handler *SignupHandler) executeFront(w http.ResponseWriter, page string, data any) error {
	err := signupFront.ExecuteTemplate(w, page, data)
	if err != nil {
		return err
	}

	return nil
}

func (handler *SignupHandler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		user := user.User{
			Name:                 strings.TrimSpace(r.FormValue("name")),
			Email:                strings.TrimSpace(r.FormValue("email")),
			Phone:                strings.TrimSpace(r.FormValue("phone")),
			CPF:                  strings.TrimSpace(r.FormValue("cpf")),
			Password:             r.FormValue("password"),
			ConfirmationPassword: r.FormValue("confirmation_password"),
		}
		validationError := handler.validator.Struct(user)
		if validationError != nil {
			dataToView := map[string]interface{}{
				"Error": validationError,
				"User":  user,
			}
			err := handler.executeFront(w, "index", dataToView)
			if err != nil {
				log.Fatal().Err(err).Msgf("canot execute template index")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		user.Password = handler.hashKit.HashField(user.Password)
		err := handler.userRepo.InsertUser(&user)
		if err != nil {
			var dataError interface{} = "Erro ao cadastrar usuário."
			dataToView := map[string]interface{}{
				"Error": dataError,
			}
			err = handler.executeFront(w, "index", dataToView)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err := signupFront.ExecuteTemplate(w, "index", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}
