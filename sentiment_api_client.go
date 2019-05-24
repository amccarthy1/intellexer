package intellexer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	listOntologiesURL    = "sentimentAnalyzerOntologies"
	analyzeSentimentsURL = "analyzeSentiments"
)

// ListOntologies lists the ontologies available for analysis
func (c *Client) ListOntologies() ([]Ontology, error) {
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
// You should assume this call will take a while. It is a network call to a
// machine learning-based API, and therefore could have a lot of overhead.
// Also, take care not to exceed the request size determined by your API level.
func (c *Client) AnalyzeSentiments(ontology Ontology, reviews []Review) (*SentimentResponse, error) {
	url := fmt.Sprintf("%s?%s", analyzeSentimentsURL, c.queryString(param{"ontology", string(ontology)}))
	res, err := c.post(url, reviews)
	if err != nil {
		return nil, err
	}
	var sentimentResponse SentimentResponse
	bytes, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(bytes))
	// decoder := json.NewDecoder(res.Body)
	// err = decoder.Decode(&sentimentResponse)
	err = json.Unmarshal(bytes, &sentimentResponse)
	if err != nil {
		return nil, err
	}
	return &sentimentResponse, nil
}
