package session

import (
	"fmt"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/cfg"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"net/http"
)

type InternalSessionService struct {
}

func NewSessionService() InternalSessionService {
	return InternalSessionService{}
}

// GenerateSessionId generates a new session id
func (service *InternalSessionService) GenerateSessionId() string {
	sessionId := fmt.Sprintf("session-id-%s", uuid.New().String())
	return sessionId
}

// CreateUserSession creates a new user session
func (service *InternalSessionService) CreateUserSession(sessionStore *sessions.CookieStore, sessionName string, sessionValue string, w http.ResponseWriter, r *http.Request) error {
	session, _ := sessionStore.Get(r, sessionName)

	session.Values[cfg.SESSION_USER_KEY] = sessionValue
	session.Values[cfg.SESSION_ID_KEY] = service.GenerateSessionId()
	session.Values[cfg.SESSION_AUTH_USER_KEY] = true
	err := session.Save(r, w)
	if err != nil {
		return err
	}

	return nil

}
