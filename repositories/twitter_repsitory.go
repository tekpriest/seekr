package repositories

import (
	"github.com/tekpriest/seekr/models"
	"github.com/tekpriest/seekr/modules/twitter"
	"go.mongodb.org/mongo-driver/mongo"
)

type Query func(models.Twitter) bool

type TwitterRepository interface {
	Save()
	Select(q Query) (twitter models.Twitter)
	InsertOrUpdate(twitter models.Twitter) (updatedTwitter models.Twitter, err error)
}

type twitterRepository struct {
	col *mongo.Collection
}

// InsertOrUpdate implements TwitterRepository.
func (t *twitterRepository) InsertOrUpdate(twitter models.Twitter) (updatedTwitter models.Twitter, err error) {
	panic("unimplemented")
}

// Select implements TwitterRepository.
func (t *twitterRepository) Select(q Query) (twitter models.Twitter) {
	panic("unimplemented")
}

func NewTwitterRepository(db twitter.DatabaseConnection) TwitterRepository {
	col := db.GetDB().DB.Collection("twitter")

	return &twitterRepository{
		col: col,
	}
}
