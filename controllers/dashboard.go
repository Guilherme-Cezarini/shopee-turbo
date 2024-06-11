package controllers

import (
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/cfg"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/middlewares"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/repository/access_token"
	user2 "github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/repository/user"
	authorization_ml "github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/service/authorization-ml"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
	"text/template"
)

var dashboardFront = template.Must(template.ParseGlob("front-end/views/dashboard/*"))

type DashboardHandler struct {
	authorizationMlService authorization_ml.AuthorizationService
	sessionStore           *sessions.CookieStore
	userRepo               user2.UserRepository
	accessTokenRepo        access_token.AccessTokenRepository
}

func NewDashboardHandler(sessionStore *sessions.CookieStore, userRepo user2.UserRepository, accessTokenRepo access_token.AccessTokenRepository) *DashboardHandler {
	return &DashboardHandler{
		authorizationMlService: authorization_ml.NewAuthorizationService(),
		sessionStore:           sessionStore,
		userRepo:               userRepo,
		accessTokenRepo:        accessTokenRepo,
	}
}

func (handler *DashboardHandler) executeFront(w http.ResponseWriter, page string, data any) error {
	err := loginFront.ExecuteTemplate(w, page, data)
	if err != nil {
		return err
	}

	return nil
}

func (handler *DashboardHandler) generateViewData(dataError any, r *http.Request) map[string]interface{} {
	data := map[string]interface{}{
		"Error":          dataError,
		csrf.TemplateTag: csrf.TemplateField(r),
	}
	return data

}

func (h *DashboardHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/dashboard/authorize-ml/callback", h.AuthorizeMl).Methods("GET")
	r.Handle("/dashboard", middlewares.AuthMiddleware(h.sessionStore)(http.HandlerFunc(h.Index))).Methods("GET")
}

// AuthorizeMl authorize the user to access the ML API
func (h *DashboardHandler) AuthorizeMl(w http.ResponseWriter, r *http.Request) {
	session, _ := h.sessionStore.Get(r, "session")
	//get query param
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	session.Values["state-ml"] = state
	userId := session.Values[cfg.SESSION_USER_KEY].(string)
	err := h.userRepo.InsertCodeInUser(userId, code)
	if err != nil {
		dataError := map[string]any{"Error": "Erro ao autorizar o acesso a API do Mercado Livre."}
		dataToView := h.generateViewData(dataError, r)
		err := h.executeFront(w, "index", dataToView)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	err, accessToken := h.authorizationMlService.RequestNewMLToken(code)
	if err != nil {
		dataError := map[string]any{"Error": "Erro ao autorizar o acesso a API do Mercado Livre."}
		dataToView := h.generateViewData(dataError, r)
		err := h.executeFront(w, "index", dataToView)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	err = h.accessTokenRepo.CreateAccessToken(userId, accessToken)
	if err != nil {
		dataError := map[string]any{"Error": "Erro ao autorizar o acesso a API do Mercado Livre."}
		dataToView := h.generateViewData(dataError, r)
		err := h.executeFront(w, "index", dataToView)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	return

}

// Index show the dashboard
func (h *DashboardHandler) Index(w http.ResponseWriter, r *http.Request) {

	err := dashboardFront.ExecuteTemplate(w, "index", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
