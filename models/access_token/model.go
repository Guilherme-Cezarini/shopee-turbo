package access_token

import "github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/models/responses"

type AccessToken struct {
	AccessTokenResponse *responses.AccessTokenResponse
	UserId              string `json:"user_id"`
}
