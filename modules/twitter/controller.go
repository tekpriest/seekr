package twitter

import (
	"context"

	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/v12"
	"github.com/robfig/cron"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	pathSaveKeyword = mvc.Response{}
)

type SaveQuery struct {
	ID        string        `bson:"id"`
	SearchKey string        `bson:"search_key"`
	results   []interface{} `bson:"results"`
}

type TwitterController interface {
	Save() string
}

type twitterController struct {
	db  *mongo.Collection
	c   *cron.Cron
	s   TwitterService
	ctx iris.Context
}

func (t *twitterController) Search() {}

func (t *twitterController) Save() string {
	var query interface{}
	if err := t.ctx.ReadQuery(&query); err != nil {
		return err.Error()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	keyword := query["keyword"]

	var search SaveQuery
	// TODO: check if keyword exists
	result := t.db.FindOne(ctx, bson.D{{Key: "searchKey", Value: keyword}})
	if err := result.Decode(&search); err != nil {
		return err.Error()
	}
	// TODO: save if not
	// TODO: create a cron job based on the keyword

	return "saved"
}

func (t *twitterController) RunQueryUpdates(keyword string) (err error) {
	t.c.AddFunc("@every 5m", func() {
	})

	return
}

func NewTwitterController(s TwitterService) TwitterController {
	c := cron.New()
	return &twitterController{
		db: &mongo.Collection{},
		c:  c,
		s:  s,
	}
}
