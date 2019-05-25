package intellexer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const (
	defaultBaseURL            = "https://api.intellexer.com"
	listOntologiesEndpoint    = "sentimentAnalyzerOntologies"
	analyzeSentimentsEndpoint = "analyzeSentiments"
)

// NewClient returns a new client with the specified API key
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
	}
}

// WithHTTPClient sets the internal HTTP client that should be used.
// Default is http.DefaultClient
func (c *Client) WithHTTPClient(client httpClient) *Client {
	c.client = client
	return c
}

// WithBaseURL sets the internal base URL to hit when sending API requests.
// Overriding is useful for testing and development.
func (c *Client) WithBaseURL(baseURL string) *Client {
	c.baseURL = baseURL
	return c
}

type httpClient interface {
	// Do sends an HTTP request
	Do(req *http.Request) (*http.Response, error)
}

// Client is an intellexer API client
type Client struct {
	baseURL string
	apiKey  string
	client  httpClient
}

// APIError is an error returned by the intellexer API. You can retrieve the response object
// from this error. These are returned wrapped as pkg/error objects for the purpose of stack
// traces, but you can get the cause by calling .Cause() on those.
type APIError struct {
	Response *http.Response
}

func (err APIError) Error() string {
	return fmt.Sprintf("Intellexer API responded with status code %d", err.Response.StatusCode)
}

type param struct {
	key   string
	value string
}

func (c *Client) queryString(params ...param) string {
	qString := url.Values(make(map[string][]string))
	qString.Add("apiKey", c.apiKey)
	for _, param := range params {
		qString.Add(param.key, param.value)
	}
	return qString.Encode()
}

func (c *Client) getPath(path string) string {
	// use custom URL if provided, otherwise default to base.
	url := c.baseURL
	if len(url) == 0 {
		url = defaultBaseURL
	}
	return fmt.Sprintf("%s/%s", url, path)
}

func (c *Client) get(path string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.getPath(path), nil)
	if err != nil {
		return nil, errors.Wrap(err, "Request creation failed")
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Request failed")
	}
	return handleResponseErrorCodes(res)
}

func (c *Client) post(path string, jsonBody interface{}) (*http.Response, error) {
	body, err := json.Marshal(jsonBody)
	if err != nil {
		return nil, errors.Wrap(err, "JSON serialization failed")
	}
	req, err := http.NewRequest("POST", c.getPath(path), bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "Request creation failed")
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Request failed")
	}
	return handleResponseErrorCodes(res)
}

func handleResponseErrorCodes(res *http.Response) (*http.Response, error) {
	if res.StatusCode >= 500 {
		err := APIError{res}
		return nil, errors.Wrap(err, "Server Error")
	}
	if res.StatusCode >= 400 {
		err := APIError{res}
		return nil, errors.Wrap(err, "Request Error")
	}
	return res, nil
}

// =============================================================================
// API Endpoint Implementations

// ListOntologies lists the ontologies available for analysis
func (c *Client) ListOntologies() ([]Ontology, error) {
	res, err := c.get(fmt.Sprintf("%s?%s", listOntologiesEndpoint, c.queryString()))
	if err != nil {
		return nil, err
	}
	var ontologies []Ontology
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ontologies)
	if err != nil {
		return nil, errors.Wrap(err, "Error deserializing response")
	}
	return ontologies, nil
}

// AnalyzeSentiments analyzes the reviews passed in for overall sentiment.
// You should assume this call will take a while. It is a network call to a
// machine learning-based API, and therefore could have a lot of overhead.
// Also, take care not to exceed the request size determined by your API level.
func (c *Client) AnalyzeSentiments(ontology Ontology, reviews []Review) (*SentimentResponse, error) {
	url := fmt.Sprintf("%s?%s", analyzeSentimentsEndpoint, c.queryString(param{"ontology", string(ontology)}))
	res, err := c.post(url, reviews)
	if err != nil {
		return nil, err
	}
	var sentimentResponse SentimentResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&sentimentResponse)
	if err != nil {
		return nil, errors.Wrap(err, "Error deserializing response")
	}
	return &sentimentResponse, nil
}
