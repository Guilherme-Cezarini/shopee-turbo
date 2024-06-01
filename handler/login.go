package handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"text/template"
)

type LoginHandler struct {
}

var loginFront = template.Must(template.ParseGlob("front-end/login/*"))

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{}
}

func (handler *LoginHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/login", handler.Index).Methods("GET")
}

func (handler *LoginHandler) Index(w http.ResponseWriter, r *http.Request) {
	err := loginFront.ExecuteTemplate(w, "index", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (handler *LoginHandler) Signup(w http.ResponseWriter, r *http.Request) {

	err := loginFront.ExecuteTemplate(w, "index", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
