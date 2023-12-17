package twitter

import (
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/v12"
	"github.com/robfig/cron"
	"go.mongodb.org/mongo-driver/mongo"
)

type SaveQuery struct {
	ID        string        `bson:"id"`
	SearchKey string        `bson:"search_key"`
	results   []interface{} `bson:"results"`
}

type TwitterController interface {
	Save() (TwitterRecentSearchResponse, string)
	Search() (TwitterRecentSearchResponse, string)
}

type twitterController struct {
	db  *mongo.Collection
	c   *cron.Cron
	s   TwitterService
	ctx iris.Context
}

func (t *twitterController) Search() (TwitterRecentSearchResponse, string) {
	var query TwitterRecentSearchQuery
	if err := t.ctx.ReadQuery(&query); err != nil {
		return TwitterRecentSearchResponse{}, err.Error()
	}
	data, err := t.s.Search(query)
	if err != nil {
		return TwitterRecentSearchResponse{}, err.Error()
	}

	return data, ""
}

func (t *twitterController) Save() (TwitterRecentSearchResponse, string) {
	var query map[string]interface{}
	if err := t.ctx.ReadQuery(&query); err != nil {
		return TwitterRecentSearchResponse{}, err.Error()
	}

	keyword := query["keyword"].(string)

	data, err := t.s.Save(keyword)
	if err != nil {
		return TwitterRecentSearchResponse{}, err.Error()
	}

	if err := t.RunQueryUpdates(keyword); err != nil {
		return data, err.Error()
	}

	return data, "saved"
}

func (t *twitterController) RunQueryUpdates(keyword string) (err error) {
	if err := t.c.AddFunc("@every 5m", func() {
		_, err := t.s.Save(keyword)
		if err != nil {
			return
		}
	}); err == nil {
		t.c.Start()
	} else {
		return err
	}

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
