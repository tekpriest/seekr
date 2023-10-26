package twitter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// RecentSearch implements TwitterService.
func (t *Twitter) RecentSearch(q TwitterRecentSearchQuery) (*TwitterRecentSearchResponse, error) {
	var data TwitterRecentSearchResponse
	query := formatQueryFields(q)
	body, err := t.doGetRequest("/tweets/search/recent", query)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(body, &data)

	return &data, nil
}

func NewTwitterService() TwitterService {
	// load all envs
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	accessToken := os.Getenv("ACCESS_TOKEN")
	tokenSecret := os.Getenv("TOKEN_SECRET")
	bearerToken := os.Getenv("BEARER_TOKEN")
	if consumerKey == "" || consumerSecret == "" || accessToken == "" || tokenSecret == "" {
		panic("Environment values not set")
	}

	twitter := &Twitter{
		BaseURL:        "https://api.twitter.com/2",
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		AccessToken:    accessToken,
		TokenSecret:    tokenSecret,
		BeaerToken:     bearerToken,
	}

	client := twitter.createHTTPClient()
	twitter.Client = *client

	return twitter
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
		break
	case "poll":
		for _, f := range AllowedPollFields {
			allowedFields[f] = true
		}
		break
	case "place":
		for _, f := range AllowedPlaceFields {
			allowedFields[f] = true
		}
		break
	case "tweet":
		for _, f := range AllowedTweetFields {
			allowedFields[f] = true
		}
		break
	case "expansion":
		for _, f := range AllowedExpansionFields {
			allowedFields[f] = true
		}
		break
	case "media":
		for _, f := range AllowedMediaFields {
			allowedFields[f] = true
		}
		break
	}

	for _, f := range fields {
		if allowedFields[f] {
			result = append(result, f)
		}
	}

	return strings.Join(result, ",")
}

func (t *Twitter) createHTTPClient() *TransportClient {
	return &TransportClient{
		BaseURL: t.BaseURL,
		Client:  &http.Client{},
	}
}

func (t *Twitter) doGetRequest(path string, q url.Values) ([]byte, error) {
	url := t.BaseURL + path

	query := q.Encode()
	if query != "" {
		url += fmt.Sprintf("?%s", q.Encode())
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}

	req.Header.Add("consumer_key", t.ConsumerKey)
	req.Header.Add("consumer_secret", t.ConsumerSecret)
	req.Header.Add("access_token", t.AccessToken)
	req.Header.Add("token_secret", t.TokenSecret)

	if err := t.preQuest(); err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}

	resp, err := t.Client.Do(req)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error creating request: %v", resp.Status)
	}
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error processing request: %v", err)
	}

	return body, nil
}

func (t *Twitter) doPostRequest(path string, q url.Values) ([]byte, error) {
	url := t.Client.BaseURL + path

	query := q.Encode()
	if query != "" {
		url += fmt.Sprintf("?%s", q.Encode())
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}

	resp, err := t.Client.Do(req)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error creating request: %v", resp.Status)
	}
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error processing request: %v", err)
	}

	return body, nil
}

func (t *Twitter) preQuest() error {
	if t.BeaerToken == "" {
		req, err := http.NewRequest(
			"POST",
			"https://api.twitter.com/oauth2/token",
			strings.NewReader("grant_type=client_credentials"),
		)

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.SetBasicAuth(t.ConsumerKey, t.ConsumerSecret)

		resp, err := t.Client.Do(req)
		if err != nil {
			return fmt.Errorf("Error logging in: %v", err)
		}
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Error processing login: %v", err)
		}
		data := map[string]interface{}{}
		if err := json.Unmarshal(b, &data); err != nil {
			return fmt.Errorf("Error parsing login details: %v", err)
		}
		if data["access_token"] == nil {
			return fmt.Errorf(
				"Error parsing login: %v",
				data["errors"].([]interface{})[0],
			)
		}

		t.BeaerToken = data["access_token"].(string)

	}

	return nil
}
