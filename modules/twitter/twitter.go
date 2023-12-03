package twitter

import "time"

type Tweet struct {
	ID        string
	TweetID string
	Text      string
	Tags      []string
	TweetedAt time.Time
	CreatedAt time.Time
}
