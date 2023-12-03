package twitter

type TwitterRepository interface {
	GetAllTweets(keyword string) ([]Tweet, error)
	SaveTweet(tweet TweetAPI) error
	GetByTag(tag string) ([]Tweet, error)
}

// type TwitterAPIRepository interface {}
