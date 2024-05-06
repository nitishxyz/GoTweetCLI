package twitter

import "net/http"

type Client struct {
	clientId     string
	clientSecret string
	httpClient   http.Client
	server       *http.Server
	oAuthToken   string
	accessToken  *AccessTokenResponse
}

func NewClient(clientId string, clientSecret string) Client {
	return Client{
		httpClient:   http.Client{},
		clientId:     clientId,
		clientSecret: clientSecret,
	}
}
