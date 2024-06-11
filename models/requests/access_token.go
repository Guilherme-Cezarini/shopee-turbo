package requests

import "github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/cfg"

type AccessToken struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	RedirectUri  string `json:"redirect_uri"`
}

// NewAccessTokenRequest creates a new AccessToken request
func NewAccessTokenRequest(code string) *AccessToken {
	return &AccessToken{
		GrantType:    "authorization_code",
		ClientID:     cfg.Env().MercadoLivreClientId,
		ClientSecret: cfg.Env().MercadoLivreClientSecret,
		Code:         code,
		RedirectUri:  cfg.Env().MercadoLivreRedirectUri,
	}
}
