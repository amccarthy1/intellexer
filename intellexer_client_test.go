package intellexer

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newMockClient(statusCode int, body string) mockClient {
	res := &http.Response{
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(body))),
		StatusCode: statusCode,
	}
	return mockClient{response: res}
}

func newMockClientFromFile(statusCode int, filename string) mockClient {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return mockClient{response: &http.Response{
		Body:       file,
		StatusCode: statusCode,
	}}
}

type mockClient struct {
	response *http.Response
	err      error
}

func (mc mockClient) Do(*http.Request) (*http.Response, error) {
	return mc.response, mc.err
}

func TestQueryString(t *testing.T) {
	client := &Client{
		apiKey: "test",
	}
	qs := client.queryString(
		param{"foo", "bar"},
		param{"foo2", "bar2"},
	)
	assert.Equal(t, "apiKey=test&foo=bar&foo2=bar2", qs)
}

func TestGetPath(t *testing.T) {
	client := &Client{
		apiKey: "test",
	}
	assert.Equal(t, "https://api.intellexer.com/foo/bar/baz", client.getPath("foo/bar/baz"))
	client.baseURL = "https://example.com"
	assert.Equal(t, "https://example.com/foo/bar/baz", client.getPath("foo/bar/baz"))
}

func TestListOntologies(t *testing.T) {
	client := newMockClient(200, `["foo", "bar", "baz"]`)
	apiClient := &Client{
		baseURL: "FAKEURL",
		apiKey:  "test",
		client:  client,
	}
	ontologies, err := apiClient.ListOntologies()
	assert.Nil(t, err)
	assert.Len(t, ontologies, 3)

	apiClient.client = newMockClientFromFile(200, "testdata/list_ontologies_response.json")
	ontologies, err = apiClient.ListOntologies()
	assert.Nil(t, err)
	assert.Len(t, ontologies, 3)
	assert.Equal(t, ontologies[0], Ontology("Hotels"))
	assert.Equal(t, ontologies[1], Ontology("Restaurants"))
	assert.Equal(t, ontologies[2], Ontology("Gadgets"))
}

func TestAnalyzeSentiments(t *testing.T) {
	client := newMockClientFromFile(200, "testdata/analyze_sentiments_response.json")
	apiClient := NewClient("test").WithBaseURL("FAKEURL").WithHTTPClient(client)
	res, err := apiClient.AnalyzeSentiments(Restaurants, NewAnalyzeSentimentsRequestBody([]string{
		"I love coffee",
		"I hate coffee",
	}))
	assert.Nil(t, err)
	assert.NotNil(t, res)
}
