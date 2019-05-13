package intellexer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const defaultBaseURL = "https://api.intellexer.com"

type httpClient interface {
	// Do sends an HTTP request
	Do(req *http.Request) (*http.Response, error)
}

type param struct {
	key   string
	value string
}

// Client is an intellexer API client
type Client struct {
	baseURL string
	apiKey  string
	client  httpClient
}

func (c Client) queryString(params ...param) string {
	qString := url.Values(make(map[string][]string))
	qString.Add("apiKey", c.apiKey)
	for _, param := range params {
		qString.Add(param.key, param.value)
	}
	return qString.Encode()
}

func (c Client) getPath(path string) string {
	// use custom URL if provided, otherwise default to base.
	url := c.baseURL
	if len(url) == 0 {
		url = defaultBaseURL
	}
	return fmt.Sprintf("%s/%s", url, path)
}

func (c Client) getHTTPClient() httpClient {
	if c.client != nil {
		return c.client
	}
	return http.DefaultClient
}

func (c Client) get(path string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.getPath(path), nil)
	if err != nil {
		return nil, errors.Wrap(err, "Request creation failed")
	}
	res, err := c.getHTTPClient().Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Request failed")
	}
	return handleResponseErrorCodes(res)
}

func (c Client) post(path string, jsonBody interface{}) (*http.Response, error) {
	body, err := json.Marshal(jsonBody)
	if err != nil {
		return nil, errors.Wrap(err, "JSON serialization failed")
	}
	req, err := http.NewRequest("POST", c.getPath(path), bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "Request creation failed")
	}
	res, err := c.getHTTPClient().Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Request failed")
	}
	return handleResponseErrorCodes(res)
}

func handleResponseErrorCodes(res *http.Response) (*http.Response, error) {
	if res.StatusCode > 500 {
		err := APIError{res}
		return nil, errors.Wrap(err, "Server Error")
	}
	if res.StatusCode > 400 {
		err := APIError{res}
		return nil, errors.Wrap(err, "Request Error")
	}
	return res, nil
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
