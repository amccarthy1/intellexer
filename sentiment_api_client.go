package intellexer

import (
	"encoding/json"
	"fmt"
)

const (
	listOntologiesURL    = "sentimentAnalyzerOntologies"
	analyzeSentimentsURL = "analyzeSentiments"
)

// ListOntologies lists the ontologies available for analysis
func (c Client) ListOntologies() ([]Ontology, error) {
	res, err := c.get(fmt.Sprintf("%s?%s", listOntologiesURL, c.queryString()))
	if err != nil {
		return nil, err
	}
	var ontologies []Ontology
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ontologies)
	return ontologies, err
}

// AnalyzeSentiments analyzes the reviews passed in for overall sentiment.
func (c Client) AnalyzeSentiments(ontology Ontology, reviews []Review) (*SentimentResponse, error) {
	url := fmt.Sprintf("%s?%s", listOntologiesURL, c.queryString(param{"ontology", string(ontology)}))
	res, err := c.post(url, reviews)
	if err != nil {
		return nil, err
	}
	var sentimentResponse SentimentResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&sentimentResponse)
	if err != nil {
		return nil, err
	}
	return &sentimentResponse, nil
}
