// Package main provides an example application making use of the intellexer API,
// and also a command-line interface for exploring the API's behavior.
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/amccarthy1/intellexer"
	"github.com/google/uuid"
)

type subcommand func(client *intellexer.Client)

var subCommands = map[string]subcommand{
	"analyze_sentiments":  analyzeSentiments,
	"get_topics":          getTopics,
	"get_topics_from_url": getTopicsFromURL,
}

func usage() {
	fmt.Println("Usage: intellexer [subcommand] (...args)")
	fmt.Println("Where subcommand is one of:")
	for option := range subCommands {
		fmt.Printf("\t%s\n", option)
	}
}

func analyzeSentiments(client *intellexer.Client) {
	if len(os.Args) < 3 {
		fmt.Println("Usage: intellexer analyze_sentiments [review] (...[reviews])")
	}
	var reviews []intellexer.Review
	for _, review := range os.Args[2:] {
		reviews = append(reviews, intellexer.Review{
			ID:   uuid.New(),
			Text: review,
		})
	}
	res, err := client.AnalyzeSentiments(intellexer.Gadgets, reviews)
	if err != nil {
		panic(err)
	}
	for _, sentiment := range res.Sentiments {
		fmt.Println(sentiment.SentimentWeight)
	}
}

func getTopics(client *intellexer.Client) {
	if len(os.Args) != 3 {
		fmt.Println("Usage: intellexer get_topics [path/to/article]")
	}
	file, err := os.Open(os.Args[2])
	if err != nil {
		panic(err)
	}
	topics, err := client.GetTopics(file)
	if err != nil {
		panic(err)
	}
	for i, topic := range topics {
		fmt.Printf("Topic %d: %s\n", i+1, topic)
	}
}

func getTopicsFromURL(client *intellexer.Client) {
	if len(os.Args) != 3 {
		fmt.Println("Usage: intellexer get_topics_from_url [url]")
	}
	topics, err := client.GetTopicsFromURL(os.Args[2])
	if err != nil {
		panic(err)
	}
	for i, topic := range topics {
		fmt.Printf("Topic %d: %s\n", i+1, topic)
	}
}

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
	subCommandStr := os.Args[1]
	subCommand, ok := subCommands[subCommandStr]
	if !ok {
		// TODO implement usage
		fmt.Printf("Invalid command: %s\n", subCommandStr)
		usage()
	}
	subCommand(client)
}
