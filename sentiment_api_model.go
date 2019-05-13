package intellexer

import (
	"github.com/google/uuid"
)

// Ontology is a context within which sentiment analysis evaluates reviews.
type Ontology string

// These are all the supported ontologies for intellexer.
// Note that the endpoint for listing these will capitalize these, but the
// sentiment analysis endpoint will not. The API is case-insensitive, so for the
// purposes of unit testing, they will be all lowercase. It is recommended to
// convert to lowercase in any code expecting equality.
const (
	Hotels      = Ontology("hotels")
	Restaurants = Ontology("restaurants")
	Gadgets     = Ontology("gadgets")
)

// Opinion is a nested set of analyzed components of a review. It contains data
// about which parts of the review contributed positively and negatively.
type Opinion struct {
	// Children is a slice of sub-opinions that are related to / factored into
	// this opinion
	Children []Opinion `json:"children"`
	// F is an undocumented field
	F int `json:"f"`
	// RS is an undocumented field
	RS []int `json:"rs"`
	// Text is either the topic of this opinion or the text from the review that
	// it is based on. This may not always come directly from the review text.
	Text *string `json:"t"`
	// SentimentWeight is the positive or negative weight of this opinion
	SentimentWeight float64 `json:"w"`
}

// Sentiment is an overall assessment of a review. A positive review will have a
// SentimentWeight above 0, and a negative one will be below 0.
type Sentiment struct {
	// ID is the unique ID of this sentiment, should be what was passed in.
	ID string `json:"id"`
	// SentimentWeight is the positive or negative weight of this review.
	SentimentWeight float64 `json:"w"`

	// These fields don't seem to ever populate and are always null.
	Author   *string `json:"author"`
	Datetime *string `json:"dt"`
	Title    *string `json:"title"`
}

// Sentence is a sentence within the review that has been annotated with an XML-
// like format. Key words/phrases are enclosed in either <pos> tags or <neg>
// tags with attribute `w` as the sentiment weight of that word.
type Sentence struct {
	// SentimentID is the ID of the sentement that this sentence comes from.
	SentimentID string `json:"sid"`
	// Text is the xml-annotated text of this sentence
	Text string `json:"text"`
	// SentimentWeight is the positive or negative weight of this sentence
	SentimentWeight float64 `json:"w"`
}

// SentimentResponse is the response format from the AnalyzeSentiments API
type SentimentResponse struct {
	SentimentsCount int         `json:"sentimentsCount"`
	Ontology        Ontology    `json:"ontology"`
	Sentences       []Sentence  `json:"sentences"`
	Opinions        Opinion     `json:"opinions"`
	Sentiments      []Sentiment `json:"sentiments"`
}

// Review is the text and ID of a review that should be analyzed for sentiment.
type Review struct {
	ID   uuid.UUID `json:"id"`
	Text string    `json:"text"`
}

// NewAnalyzeSentimentsRequestBody returns a new request body for the
// /analyzeSentiments endpoint, generating UUIDs for each review. It is not
// recommended to use this, callers are instead recommended to generate their
// own UUIDs so they can be cross-referenced with the results.
func NewAnalyzeSentimentsRequestBody(reviews []string) []Review {
	var sentimentRequests []Review
	for _, review := range reviews {
		sentimentRequests = append(sentimentRequests, Review{
			ID:   uuid.New(),
			Text: review,
		})
	}
	return sentimentRequests
}
