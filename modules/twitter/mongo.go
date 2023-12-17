package twitter

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/tekpriest/seekr/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	DB *mongo.Database
}

func NewDatabaseConnection(c *config.Config) (*MongoRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 15)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(c.DBURL))
	if err != nil {
		panic(errors.Wrap(err, "there was an error connecting to the db"))
	}
	db := client.Database("seekr")

	return &MongoRepository{
	DB: db,
}, nil
}
