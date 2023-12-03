package twitter

import "time"

type TwitterRecentSearchQuery struct {
	Query           string    `json:"query,omitempty" url:"query"`
	TweetFields     []string  `json:"tweet_fields,omitempty" url:"tweet_fields"`
	StartTime       time.Time `json:"start_time,omitempty" url:"start_time"`
	EndTime         time.Time `json:"end_time,omitempty" url:"end_time"`
	SinceID         string    `json:"since_id,omitempty" url:"since_id"`
	UntilID         string    `json:"until_id,omitempty" url:"until_id"`
	MaxResults      string    `json:"max_results,omitempty" url:"max_results"`
	NextToken       string    `json:"next_token,omitempty" url:"next_token"`
	ExpansionFields []string  `json:"expansions,omitempty" url:"expansion_fields"`
	PlaceFields     []string  `json:"place_fields,omitempty" url:"place_fields"`
	PollFields      []string  `json:"poll_fields,omitempty" url:"poll_fields"`
	UserFields      []string  `json:"user_fields,omitempty" url:"user_fields"`
}
