// Package mocks provides utilities for testing the intellexer client. Mostly
// this includes HTTP client mocks that can simulate simple success and error
// cases.
package mocks

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// NewMockClient returns a new mock client that always responds as instructed
func NewMockClient(statusCode int, body string) MockClient {
	return MockClient{
		body: []byte(body),
		statusCode: statusCode,
	}
}

// NewMockClientFromFile mocks an HTTP client that responds with the contents
// of a file.
func NewMockClientFromFile(statusCode int, filename string) MockClient {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return MockClient{
		statusCode: statusCode,
		body: body,
	}
}

// NewErrorClient returns a client that always errors on requests.
func NewErrorClient(err error) MockClient {
	return MockClient{err: err}
}

// MockClient is a fake HTTP client that responds with a static response or
// error. It is useful for unit testing the client itself.
type MockClient struct {
	statusCode int
	body       []byte
	err        error
}

// Do fakes an HTTP response without actually sending a request.
func (mc MockClient) Do(*http.Request) (*http.Response, error) {
	if mc.err != nil {
		return nil, mc.err
	}
	return &http.Response{
		Body:       ioutil.NopCloser(bytes.NewReader(mc.body)),
		StatusCode: mc.statusCode,
	}, nil
}
