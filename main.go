package main

import (
	"fmt"
	"os"

	"github.com/tekpriest/seekr/pkg/twitter"
)

func main() {
	t := twitter.NewTwitterService()

	q := twitter.TwitterRecentSearchQuery{
		Query:       "doggo",
		TweetFields: []string{"attachments"},
	}

	tweets, err := t.RecentSearch(q)
	if err != nil {
    fmt.Println(err)
    os.Exit(1)
	}
	fmt.Println(tweets)
}
