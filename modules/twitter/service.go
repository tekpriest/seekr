package twitter

import (
	"github.com/tekpriest/seekr/config"
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
