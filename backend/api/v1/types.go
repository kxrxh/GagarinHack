package v1

type RequestBody struct {
	RequestMessage string `json:"request_message"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}
