package authorize

import (
	"fmt"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/cfg"
	"github.com/google/uuid"
)

type authorizeService struct {
}

type AuthorizeService interface {
}

func NewService() AuthorizeService {
	return &authorizeService{}
}

// GenerateUrlToMLAuthorization generates the URL to redirect the user to the Mercado Livre authorization page
func (service *authorizeService) GenerateUrlToMLAuthorization() string {
	stateUuid := uuid.New().String()
	redirectUri := fmt.Sprintf("%s/api/v1/authorization/callback", cfg.Env().MercadoLivreRedirectUri)
	return fmt.Sprintf("%s/authorization?response_type=code&client_id=%s&redirect_uri=%s&state=%s", cfg.Env().MercadoLivreAuthUrl, cfg.Env().MercadoLivreClientId, redirectUri, string(stateUuid))
}
