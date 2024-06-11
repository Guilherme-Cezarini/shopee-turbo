package authorization_ml

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/models/requests"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/models/responses"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"

	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/cfg"
	"github.com/google/uuid"
)

const callbackAuthUrl = "dashboard/authorize-ml/callback"

type authorizationService struct {
}

type AuthorizationService interface {
	GenerateUrlToMLAuthorization() string
	RequestNewMLToken(code string) (error, *responses.AccessTokenResponse)
}

func NewAuthorizationService() AuthorizationService {
	return &authorizationService{}
}

// GenerateUrlToMLAuthorization generates the URL to redirect the user to the Mercado Livre authorization page
func (service *authorizationService) GenerateUrlToMLAuthorization() string {
	stateUuid := uuid.New().String()
	redirectUri := fmt.Sprintf("%s/%s", cfg.Env().MercadoLivreRedirectUri, callbackAuthUrl)
	return fmt.Sprintf("%s/authorization?response_type=code&client_id=%s&redirect_uri=%s&state=%s", cfg.Env().MercadoLivreAuthUrl, cfg.Env().MercadoLivreClientId, redirectUri, string(stateUuid))
}

// RequestNewMLToken requests a new token to Mercado Livre
func (service *authorizationService) RequestNewMLToken(code string) (error, *responses.AccessTokenResponse) {
	accessTokenRequest := requests.NewAccessTokenRequest(code)
	bodyData, _ := json.Marshal(accessTokenRequest)
	req, err := http.NewRequest("POST", cfg.Env().MercadoLivreTokenUrl, bytes.NewBuffer(bodyData))
	if err != nil {
		log.Debug().Err(err).Msg("Error creating request")
		return err, nil
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Debug().Err(err).Msg("Error making request")
		return err, nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Debug().Err(err).Msg("Error reading body")
		return err, nil

	}
	data := &responses.AccessTokenResponse{}
	err = json.Unmarshal(body, data)
	if err != nil {
		log.Debug().Err(err).Msg("Error unmarshalling body")
		return err, nil
	}

	return nil, data
}
