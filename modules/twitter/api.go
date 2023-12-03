package twitter

import (
	"net/http"

	"github.com/tekpriest/seekr/config"
)

type TwitterAPIRepository struct {
	c *http.Client
}

type TransportClient struct {
	BaseURL string
	Headers map[string]string
	*http.Client
}

type TwitterAPI struct {
	ConsumerKey    string `json:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret"`
	AccessToken    string `json:"access_token"`
	TokenSecret    string `json:"token_secret"`
	BeaerToken     string
	BaseURL        string
}

func NewTwitterAPIService(c *config.Config) *TwitterAPIRepository {
	twitter := &TwitterAPI{
		BaseURL:        "https://api.twitter.com/2",
		ConsumerKey:    c.TwitterCFG.ConsumerKey,
		ConsumerSecret: c.TwitterCFG.ConsumerSecret,
		AccessToken:    c.TwitterCFG.AccessToken,
		TokenSecret:    c.TwitterCFG.TokenSecret,
		BeaerToken:     c.TwitterCFG.BearerToken,
	}

	client := createHTTPClient(c.TwitterCFG.BaseURL)

	return &TwitterAPIRepository{
		c: &client,
	}
}

func createHTTPClient(url string) *TransportClient {
	return &TransportClient{
		BaseURL: url,
		Client:  &http.Client{},
	}
}
