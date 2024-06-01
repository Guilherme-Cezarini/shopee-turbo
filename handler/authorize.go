package handler

import (
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/service/authorize"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
)

type AuthorizeHandler struct {
	authService authorize.AuthorizeService
}

func NewAuthorizeHandler() *AuthorizeHandler {
	return &AuthorizeHandler{
		authService: authorize.NewService(),
	}

}

func (h *AuthorizeHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/authorize", h.Authorize)
}

func (h *AuthorizeHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	url := h.authService.GenerateUrlToMLAuthorization()
	log.Debug().Msgf("Redirecting to %s", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
