package responses

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	UserId       int    `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
}
