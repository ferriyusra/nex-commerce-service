package response

type SuccessAuthResponse struct {
	Meta
	AccessToken string `json:"accessToken"`
	ExpiresAt   int64  `json:"expiresAt"`
}
