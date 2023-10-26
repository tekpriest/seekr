package twitter

import (
	"net/http"
	"time"
)

type ReferrencedTweets struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type EntityTimeFrame struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type TwitterMention struct {
	EntityTimeFrame
	Username string `json:"username"`
}

type Hashtag struct {
	EntityTimeFrame
	Tag string `json:"tag"`
}

type Image struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type URL struct {
	EntityTimeFrame
	URL         string  `json:"url"`
	ExpandedURL string  `json:"expanded_url"`
	DisplayURL  string  `json:"display_url"`
	Images      []Image `json:"images"`
	Status      int     `json:"status"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	UnwoundURL  string  `json:"unwound_url"`
}

type TwitterEntities struct {
	Mentions []TwitterMention `json:"mentions"`
	Hashtags []Hashtag        `json:"hashtags"`
	URLS     []URL            `json:"urls"`
}

type MediaKey string

type Attachment struct {
	MediaKeys []MediaKey `json:"media_keys"`
}

type PublicMetric struct {
	ViewCount int `json:"view_count"`
}

type IncludeMedia struct {
	Height          int          `json:"height"`
	MediaKey        MediaKey     `json:"media_key"`
	Type            string       `json:"type"`
	URL             string       `json:"url"`
	Width           int          `json:"width"`
	PreviewImageURL string       `json:"preview_image_url"`
	PublicMetrics   PublicMetric `json:"public_metrics"`
	DurationMS      int          `json:"duration_ms"`
}

type IncludeUser struct {
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Username    string    `json:"username"`
}

type Tweet struct {
	AuthorID          string              `json:"author_id"`
	CreatedAt         time.Time           `json:"created_at"`
	Entities          TwitterEntities     `json:"entities"`
	ID                string              `json:"id"`
	Lang              string              `json:"lang"`
	ReferrencedTweets []ReferrencedTweets `json:"referrenced_tweets"`
	Source            string              `json:"source"`
	Text              string              `json:"text"`
}

type (
	IncludeTweet Tweet
	Include      struct {
		Media  []IncludeMedia `json:"media"`
		User   []IncludeUser  `json:"user"`
		Tweets []IncludeTweet `json:"tweets"`
	}
)

type Geo struct {
	PlaceID string `json:"place_id"`
}

type TwitterRecentSearchResponseData struct {
	Tweet
	Attachments       []Attachment `json:"attachments"`
	PossiblySensitive bool         `json:"possibly_sensitive"`
	InReplyToUserID   string       `json:"in_reply_to_user_id"`
	Geo               Geo          `json:"geo"`
}

type ResponseMeta struct {
	NewestID    string
	OldestID    string
	ResultCount int
	NextToken   string
}

type TwitterRecentSearchResponse struct {
	Data     TwitterRecentSearchResponseData `json:"data"`
	Includes []Include                       `json:"includes"`
	Meta     ResponseMeta                    `json:"meta"`
}

type TwitterService interface {
	RecentSearch(q TwitterRecentSearchQuery) (*TwitterRecentSearchResponse, error)
}

type TwitterRecentSearchQuery struct {
	Query           string    `json:"query"`
	TweetFields     []string  `json:"tweet_fields"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	SinceID         string    `json:"since_id"`
	UntilID         string    `json:"until_id"`
	MaxResults      string    `json:"max_results"`
	NextToken       string    `json:"next_token"`
	ExpansionFields []string  `json:"expansions"`
	PlaceFields     []string  `json:"place_fields"`
	PollFields      []string  `json:"poll_fields"`
	UserFields      []string  `json:"user_fields"`
}

var AllowedTweetFields = []string{
	"attachments",
	"author_id",
	"context_annotations",
	"created_at",
	"entities",
	"geo",
	"id",
	"in_reply_to_user_id",
	"lang",
	"non_public_metrics",
	"organic_metrics",
	"possibly_sensitive",
	"promoted_metrics",
	"public_metrics",
	"referenced_tweets",
	"source",
	"text",
	"withheld",
}

var AllowedExpansionFields = []string{
	"attachments.poll_ids",
	"attachments.media_keys",
	"author_id",
	"geo.place_id",
	"in_reply_to_user_id",
	"referenced_tweets.id",
}

var AllowedMediaFields = []string{
	"duration_ms",
	"height",
	"media_key",
	"non_public_metrics",
	"organic_metrics",
	"preview_image_url",
	"promoted_metrics",
	"public_metrics",
	"type",
	"url",
	"width",
}

var AllowedPlaceFields = []string{
	"contained_within",
	"country",
	"country_code",
	"full_name",
	"geo",
	"id",
	"name",
	"place_type",
}

var AllowedPollFields = []string{
	"duration_minutes",
	"end_datetime",
	"id",
	"options",
	"voting_status",
}

var AllowedUserFields = []string{
	"created_at",
	"description",
	"entities",
	"id",
	"location",
	"name",
	"pinned_tweet_id",
	"profile_image_url",
	"protected",
	"public_metrics",
	"url",
	"username",
	"verified",
	"withheld",
}

type TransportClient struct {
	BaseURL string
	Headers map[string]string
	*http.Client
}

type Twitter struct {
	ConsumerKey    string `json:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret"`
	AccessToken    string `json:"access_token"`
	TokenSecret    string `json:"token_secret"`
	BeaerToken     string
	BaseURL        string
	Client         TransportClient
}
