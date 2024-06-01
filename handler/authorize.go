package handler

import "net/http"

type AuthorizeHandler struct {
}

func NewAuthorizeHandler() *AuthorizeHandler {
	return &AuthorizeHandler{}

}

func (h *AuthorizeHandler) RegisterRoutes() {

	http.HandleFunc("/authorize", h.Authorize)
}

func (h *AuthorizeHandler) Authorize(w http.ResponseWriter, r *http.Request) {

}
