package mocks

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
)

// NewMockClient returns a new mock client that always responds as instructed
func NewMockClient(statusCode int, body string) MockClient {
	res := &http.Response{
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(body))),
		StatusCode: statusCode,
	}
	return MockClient{response: res}
}

// NewMockClientFromFile mocks an HTTP client that responds with the contents
// of a file.
func NewMockClientFromFile(statusCode int, filename string) MockClient {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return MockClient{response: &http.Response{
		Body:       file,
		StatusCode: statusCode,
	}}
}

// NewErrorClient returns a client that always errors on requests.
func NewErrorClient(err error) MockClient {
	return MockClient{err: err}
}

// MockClient is a fake HTTP client that responds with a static response or
// error. It is useful for unit testing the client itself.
type MockClient struct {
	response *http.Response
	err      error
}

// Do fakes an HTTP response without actually sending a request.
func (mc MockClient) Do(*http.Request) (*http.Response, error) {
	return mc.response, mc.err
}
