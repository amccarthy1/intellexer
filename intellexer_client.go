// Package intellexer provides an API client implementation for various
// endpoints in the Intellexer Natural Language Processing API.
package intellexer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

const (
	defaultBaseURL            = "https://api.intellexer.com"
	getTopicsFromURLEndpoint  = "getTopicsFromUrl"
	getTopicsFromFileEndpoint = "getTopicsFromFile"
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

func (c *Client) decodeRes(res *http.Response, out interface{}) error {
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(out); err != nil {
		return errors.Wrap(err, "Error deserializing response")
	}
	return nil
}

func (c *Client) getJSON(path string, out interface{}) error {
	res, err := c.get(path)
	if err != nil {
		return err
	}
	return c.decodeRes(res, out)
}

func (c *Client) postJSON(path string, jsonBody, out interface{}) error {
	res, err := c.post(path, jsonBody)
	if err != nil {
		return err
	}
	return c.decodeRes(res, out)
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

// GetTopicsFromURL gets a list of topics from the article at the given URL.
// See doc for "GetTopics" for performance information.
func (c *Client) GetTopicsFromURL(url string) ([]string, error) {
	var topics []string
	return topics, c.getJSON(
		fmt.Sprintf("%s?%s", getTopicsFromURLEndpoint, c.queryString(param{"url", url})),
		&topics,
	)
}

// GetTopics gets a list of topics from the article read from the body.
// Note that this will actually cause the remote server to read through and
// analyze the entire article, which will usually take a few seconds and tends
// to scale with the size of the article.
func (c *Client) GetTopics(body io.Reader) ([]string, error) {
	var topics []string
	url := fmt.Sprintf("%s?%s", getTopicsFromFileEndpoint, c.queryString())
	req, err := http.NewRequest("POST", c.getPath(url), body)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating request")
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error sending request")
	}
	res, err = handleResponseErrorCodes(res)
	if err != nil {
		return nil, err
	}
	return topics, c.decodeRes(res, &topics)
}

// GetTopicsFromText is a convenience function to get topics from a string.
// You probably want to use GetTopics instead if you already have an io.Reader.
func (c *Client) GetTopicsFromText(body string) ([]string, error) {
	reader := strings.NewReader(body)
	return c.GetTopics(reader)
}

// ListOntologies lists the ontologies available for analysis. This endpoint is
// supported almost exclusively for completeness, the intellexer API only
// supports three ontologies, 'Hotels', 'Restaurants', and 'Gadgets' which are
// exported as `Hotels`, `Restaurants` and `Gadgets`.
func (c *Client) ListOntologies() ([]Ontology, error) {
	var ontologies []Ontology
	return ontologies, c.getJSON(
		fmt.Sprintf("%s?%s", listOntologiesEndpoint, c.queryString()),
		&ontologies,
	)
}

// AnalyzeSentiments analyzes the reviews passed in for overall sentiment.
// You should assume this call will take a while. It is a network call to a
// machine learning-based API, and therefore could have a lot of overhead.
// Also, take care not to exceed the request size determined by your API level.
func (c *Client) AnalyzeSentiments(ontology Ontology, reviews []Review) (*SentimentResponse, error) {
	url := fmt.Sprintf("%s?%s", analyzeSentimentsEndpoint, c.queryString(param{"ontology", string(ontology)}))
	var sentimentResponse SentimentResponse
	if err := c.postJSON(url, reviews, &sentimentResponse); err != nil {
		return nil, err
	}
	return &sentimentResponse, nil
}
