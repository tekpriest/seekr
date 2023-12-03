package twitter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/tekpriest/seekr/config"
)

type TwitterAPIRepository struct {
	c *TransportClient
	tc config.TwitterConfig
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
	client := createHTTPClient(c.TwitterCFG.BaseURL)

	return &TwitterAPIRepository{
		c: client,
	}
}

func createHTTPClient(url string) *TransportClient {
	return &TransportClient{
		BaseURL: url,
		Client:  &http.Client{},
	}
}

func (t *TwitterAPIRepository) RecentSearch(q TwitterRecentSearchQuery) (*TwitterRecentSearchResponse, error) {
	var data TwitterRecentSearchResponse
	query := formatQueryFields(q)
	body, err := t.doGetRequest("/tweets/search/recent", query)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func formatQueryFields(q TwitterRecentSearchQuery) url.Values {
	queries := url.Values{}

	if len(q.TweetFields) > 0 {
		queries.Add("tweet.fields", formatFields(q.TweetFields, "tweet"))
	}
	if len(q.ExpansionFields) > 0 {
		queries.Add("expansions", formatFields(q.ExpansionFields, "expansion"))
	}
	if len(q.PlaceFields) > 0 {
		queries.Add("place.fields", formatFields(q.PlaceFields, "place"))
	}
	if len(q.PollFields) > 0 {
		queries.Add("poll.fields", formatFields(q.PlaceFields, "poll"))
	}
	if len(q.UserFields) > 0 {
		queries.Add("user.fields", formatFields(q.PlaceFields, "user"))
	}
	if q.Query != "" {
		queries.Add("query", q.Query)
	}

	return queries
}

func formatFields(fields []string, field string) string {
	allowedFields := make(map[string]bool)
	result := []string{}

	switch field {
	case "user":
		for _, f := range AllowedUserFields {
			allowedFields[f] = true
		}
	case "poll":
		for _, f := range AllowedPollFields {
			allowedFields[f] = true
		}
	case "place":
		for _, f := range AllowedPlaceFields {
			allowedFields[f] = true
		}
	case "tweet":
		for _, f := range AllowedTweetFields {
			allowedFields[f] = true
		}
	case "expansion":
		for _, f := range AllowedExpansionFields {
			allowedFields[f] = true
		}
	case "media":
		for _, f := range AllowedMediaFields {
			allowedFields[f] = true
		}
	}

	for _, f := range fields {
		if allowedFields[f] {
			result = append(result, f)
		}
	}

	return strings.Join(result, ",")
}

func (t *TwitterAPIRepository) doGetRequest(path string, q url.Values) ([]byte, error) {
	url := t.c.BaseURL + path

	query := q.Encode()
	if query != "" {
		url += fmt.Sprintf("?%s", q.Encode())
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("consumer_key", t.tc.ConsumerKey)
	req.Header.Add("consumer_secret", t.tc.ConsumerSecret)
	req.Header.Add("access_token", t.tc.AccessToken)
	req.Header.Add("token_secret", t.tc.TokenSecret)

	if err = t.preQuest(); err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := t.c.Do(req)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error creating request: %v", resp.Status)
	}
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error processing request: %v", err)
	}

	return body, nil
}

func (t *TwitterAPIRepository) addHeaders() { }

func (t *TwitterAPIRepository) doPostRequest(path string, q url.Values) ([]byte, error) {
	url := t.c.BaseURL + path

	query := q.Encode()
	if query != "" {
		url += fmt.Sprintf("?%s", q.Encode())
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := t.c.Do(req)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error creating request: %v", resp.Status)
	}
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error processing request: %v", err)
	}

	return body, nil
}

func (t *TwitterAPIRepository) preQuest() error {
	if t.tc.BearerToken == "" {
		req, _ := http.NewRequest(
			"POST",
			"https://api.twitter.com/oauth2/token",
			strings.NewReader("grant_type=client_credentials"),
		)

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.SetBasicAuth(t.tc.ConsumerKey, t.tc.ConsumerSecret)

		resp, err := t.c.Do(req)
		if err != nil {
			return fmt.Errorf("error logging in: %v", err)
		}
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error processing login: %v", err)
		}
		data := map[string]interface{}{}
		if err := json.Unmarshal(b, &data); err != nil {
			return fmt.Errorf("error parsing login details: %v", err)
		}
		if data["access_token"] == nil {
			return fmt.Errorf(
				"error parsing login: %v",
				data["errors"].([]interface{})[0],
			)
		}

		t.tc.BearerToken = data["access_token"].(string)

	}

	return nil
}
