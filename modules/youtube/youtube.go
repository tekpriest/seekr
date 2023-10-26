package youtube

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/option"
	yt "google.golang.org/api/youtube/v3"
)

var ctx = context.TODO()

const API_KEY = "AIzaSyCtoTlIkegS-zRlSFvIZou81x5fpY16-ME"

func main() {
	service, err := yt.NewService(ctx, option.WithAPIKey(API_KEY))
	handleError(err, "Error creating Youtube client")

	searchVideosByQuery(service, []string{"snippet"}, "among us")
}

func searchVideosByQuery(s *yt.Service, parts []string, query string) {
	call := s.Search.List(parts)
	call = call.Q(query)
	response, err := call.Do()
	handleError(err, "")
	fmt.Println(
		fmt.Printf(
			"search resulted in %d results, first title: %s",
			response.PageInfo.TotalResults,
			response.Items[0].Snippet.Title,
		),
	)
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatal(fmt.Sprintf("%s: %v", message, err.Error()))
	}
}
