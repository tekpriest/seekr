package main

import (
	"time"

	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/tekpriest/seekr/config"
	"github.com/tekpriest/seekr/modules/twitter"
)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")

	// api := app.Party("/api")
	// {
	// 	api.UseRouter()
	// }

	twitterRouter := app.Party("/twitter")
	{
		manager := sessions.New(sessions.Config{
			Cookie:  "sessioncookiename",
			Expires: 24 * time.Hour,
		})
		twitterRouter.Use(manager.Handler())
		mvc.Configure(twitterRouter, configureTwitterMVC)
	}
	// twitterRouter.Handle(new(twitter.TwitterController))
	// mvc.Configure(party router.Party, configurators ...func(*mvc.Application))
	// twitterAPI
	// {
	// 	twitterAPI.Use(iris.Compression)
	// 	twitterAPI.Get("/search", search)
	// }

	if err := app.Listen(":8080", defaults); err != nil {
		panic(err)
	}
}

func defaults(app *iris.Application) {
	app.Configure(
		iris.WithOptimizations,
		iris.WithFireMethodNotAllowed,
		iris.WithLowercaseRouting,
		iris.WithPathIntelligence,
		iris.WithTunneling,
	)
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("Message", ctx.Values().GetStringDefault("message", "The page you're looking for doesn't exist"))
		if err := ctx.View("shared/error.html"); err != nil {
			ctx.HTML("<h3>%s</h3>", err.Error())
			return
		}
	})
}

func configureTwitterMVC(twitterApp *mvc.Application) {
	c := config.NewConfig()
	twitterApp.Register(twitter.NewDatabaseConnection(c))
	twitterApp.Handle(new(twitter.TwitterController))
}

// func search(ctx iris.Context) {
// 	var query twitter.TwitterRecentSearchQuery
// 	if err := ctx.ReadQuery(&query); err != nil {
// 		ctx.StopWithError(iris.StatusBadRequest, err)
// 	}
// 	ctx.JSON(query)
// }
//
// func init() {
// 	twitter.NewTwitterService()
// }
