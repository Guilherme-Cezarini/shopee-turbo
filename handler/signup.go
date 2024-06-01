package handler

import (
	"fmt"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/models/user"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
	"text/template"
)

type SignupHandler struct {
	validator *validator.Validate
}

var signupFront = template.Must(template.ParseGlob("front-end/signup/*"))

func NewSignupHandler() *SignupHandler {
	return &SignupHandler{
		validator: validator.New(),
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

func (handler *SignupHandler) Signup(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		user := user.User{
			Name:     r.FormValue("name"),
			Email:    r.FormValue("email"),
			Phone:    r.FormValue("phone"),
			CPF:      r.FormValue("cpf"),
			Password: r.FormValue("password"),
		}
		validationError := handler.validator.Struct(user)
		if validationError != nil {
			fmt.Println(validationError)
			err := signupFront.ExecuteTemplate(w, "index", validationError)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	err := signupFront.ExecuteTemplate(w, "index", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
