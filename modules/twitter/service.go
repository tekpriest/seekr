package twitter

import (
	"context"

	"github.com/lucsky/cuid"
	"github.com/tekpriest/seekr/config"
	"go.mongodb.org/mongo-driver/bson"
)

type TwitterService struct {
	MR  MongoRepository
	TAR TwitterAPIRepository
}

type TwitterServiceConfiguration func(ts *TwitterService) error

func NewService(cfgs ...TwitterServiceConfiguration) (*TwitterService, error) {
	ts := &TwitterService{}

	for _, cfg := range cfgs {
		if err := cfg(ts); err != nil {
			return nil, err
		}
	}

	return ts, nil
}

func RegisterMongoRepository(c *config.Config) TwitterServiceConfiguration {
	return func(ts *TwitterService) error {
		mr, err := NewDatabaseConnection(c)
		if err != nil {
			return err
		}
		ts.MR = *mr
		return nil
	}
}

func RegisterTwitterAPIService(c *config.Config) TwitterServiceConfiguration {
	return func(ts *TwitterService) error {
		tar := NewTwitterAPIService(c)
		ts.TAR = *tar

		return nil
	}
}

func (t *TwitterService) Search(q TwitterRecentSearchQuery) (TwitterRecentSearchResponse, error) {
	response, err := t.TAR.RecentSearch(q)
	if err != nil {
		return TwitterRecentSearchResponse{}, err
	}

	return *response, nil
}

func (t *TwitterService) Save(keyword string) (TwitterRecentSearchResponse, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var search SaveQuery
	result := t.MR.DB.Collection("tweets").FindOne(ctx, bson.D{{Key: "searchKey", Value: keyword}})
	if err := result.Decode(&search); err != nil {
		return TwitterRecentSearchResponse{}, err
	}

	q := TwitterRecentSearchQuery{
		Query: keyword,
	}
	// TODO: save if not
	// TODO: create a cron job based on the keyword
	response, err := t.Search(q)
	if err != nil {
		return TwitterRecentSearchResponse{}, err
	}
	newTweet := Tweet{ID: cuid.New()}
	_, err = t.MR.DB.Collection("tweets").InsertOne(ctx, newTweet)
	if err != nil {
		return TwitterRecentSearchResponse{}, err
	}
	return response, nil
}
