package intellexer

import (
	"github.com/google/uuid"
)

type ontology string

// These are all the supported ontologies for intellexer. An intellexer ontology
// is the target domain of a review.
const (
	Hotels      = ontology("Hotels")
	Restaurants = ontology("Restaurants")
	Gadgets     = ontology("Gadgets")
)

// Opinion is a nested set of analyzed components of a review. It contains data
// about which parts of the review contributed positively and negatively.
type Opinion struct {
	// Children is a slice of sub-opinions that are related to / factored into
	// this opinion
	Children []Opinion `json:"children"`
	// F is an undocumented field
	F int64 `json:"f"`
	// RS is an undocumented field
	RS []int64 `json:"rs"`
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

type sentimentResponse struct {
	SentimentsCount int        `json:"sentimentsCount"`
	Ontology        ontology   `json:"ontology"`
	Sentences       []Sentence `json:"sentences"`
	Opinions        []Opinion  `json:"opinions"`
}

type sentimentRequest struct {
	ID   uuid.UUID `json:"id"`
	Text string    `json:"text"`
}

// AnalyzeSentimentsRequest is a request to the /analyzeSentiments endpoint
type AnalyzeSentimentsRequest []sentimentRequest

// NewAnalyzeSentimentsRequestBody returns a new request body for the
// /analyzeSentiments endpoint, generating UUIDs for each review. It is not
// recommended to use this, callers are instead recommended to generate their
// own UUIDs so they can be cross-referenced with the results.
func NewAnalyzeSentimentsRequestBody(reviews []string) AnalyzeSentimentsRequest {
	var sentimentRequests AnalyzeSentimentsRequest
	for _, review := range reviews {
		sentimentRequests = append(sentimentRequests, sentimentRequest{
			ID:   uuid.New(),
			Text: review,
		})
	}
	return sentimentRequests
}
