Intellexer Client for Golang
==========
An API library for accessing the intellexer sentiment analysis API.

This currently does not support any APIs besides the sentiment analysis endpoints.

## Installation
`go get github.com/amccarthy1/intellexer`

This package supports go modules!

## Usage
```go
client := intellexer.NewClient(apiKey).WithHTTPClient(http.DefaultClient)
    review := intellexer.Review{
        ID:   uuid.New(),
        Text: "This gadget is neat",
    }
    res, err := client.AnalyzeSentiments(
        intellexer.Gadgets,
        []intellexer.Review{review}
    )
```
