package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/amccarthy1/intellexer"
	"github.com/google/uuid"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No arguments given")
		os.Exit(1)
	}
	apiKey, found := os.LookupEnv("API_KEY")
	if !found {
		panic("No API key given! Please run with API_KEY={api_key_here}")
	}
	client := intellexer.NewClient(apiKey).WithHTTPClient(http.DefaultClient)
	review := intellexer.Review{
		ID:   uuid.New(),
		Text: os.Args[1],
	}
	res, err := client.AnalyzeSentiments(intellexer.Gadgets, []intellexer.Review{review})
	if err != nil {
		panic(err)
	}
	fmt.Println(res.Sentiments[0].SentimentWeight)
}
