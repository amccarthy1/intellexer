package intellexer

import (
	"strings"
	"testing"

	"github.com/amccarthy1/intellexer/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

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

func TestGetTopics(t *testing.T) {
	article := "I'm an article about tech health care"
	reader := strings.NewReader(article)
	client := mocks.NewMockClientFromFile(200, "testdata/get_topics_response.json")
	apiClient := NewClient("test").WithHTTPClient(client).WithBaseURL("FAKEURL")
	topics, err := apiClient.GetTopics(reader)
	assert.Nil(t, err)
	assert.Len(t, topics, 2)
	assert.Equal(t, topics[0], "Health.healthcare")
	assert.Equal(t, topics[1], "Tech.information_technology")

	apiClient.WithHTTPClient(mocks.NewMockClientFromFile(200, "testdata/get_topics_response.json"))
	topics, err = apiClient.GetTopicsFromText(article)
	assert.Nil(t, err)
	assert.Len(t, topics, 2)
	assert.Equal(t, topics[0], "Health.healthcare")
	assert.Equal(t, topics[1], "Tech.information_technology")

	apiClient.WithHTTPClient(mocks.NewMockClientFromFile(200, "testdata/get_topics_response.json"))
	topics, err = apiClient.GetTopicsFromURL("blah/article_about_tech.php")
	assert.Nil(t, err)
	assert.Len(t, topics, 2)
	assert.Equal(t, topics[0], "Health.healthcare")
	assert.Equal(t, topics[1], "Tech.information_technology")
}

func TestListOntologies(t *testing.T) {
	client := mocks.NewMockClient(200, `["foo", "bar", "baz"]`)
	apiClient := &Client{
		baseURL: "FAKEURL",
		apiKey:  "test",
		client:  client,
	}
	ontologies, err := apiClient.ListOntologies()
	assert.Nil(t, err)
	assert.Len(t, ontologies, 3)

	apiClient.client = mocks.NewMockClientFromFile(200, "testdata/list_ontologies_response.json")
	ontologies, err = apiClient.ListOntologies()
	assert.Nil(t, err)
	assert.Len(t, ontologies, 3)
	assert.Equal(t, ontologies[0], Ontology("Hotels"))
	assert.Equal(t, ontologies[1], Ontology("Restaurants"))
	assert.Equal(t, ontologies[2], Ontology("Gadgets"))
}

func TestAnalyzeSentiments(t *testing.T) {
	client := mocks.NewMockClientFromFile(200, "testdata/analyze_sentiments_response.json")
	apiClient := NewClient("test").WithBaseURL("FAKEURL").WithHTTPClient(client)
	body := NewAnalyzeSentimentsRequestBody([]string{"I love coffee", "I hate coffee"})
	res, err := apiClient.AnalyzeSentiments(Restaurants, body)
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestAPIErrors(t *testing.T) {
	client := mocks.NewMockClientFromFile(400, "testdata/content_type_error.xhtml")
	apiClient := NewClient("test").WithBaseURL("FAKEURL").WithHTTPClient(client)
	body := NewAnalyzeSentimentsRequestBody([]string{"foo", "bar"})
	res, err := apiClient.AnalyzeSentiments(Restaurants, body)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Request Error")
	// Make sure this is an API error wrapped as pkg.error
	cause := errors.Cause(err)
	apiError, ok := cause.(APIError)
	assert.True(t, ok)
	assert.Equal(t, "Intellexer API responded with status code 400", apiError.Error())
	assert.NotNil(t, apiError.Response)
}

func TestDeserializationErrors(t *testing.T) {
	client := mocks.NewMockClient(200, "<html>I'm an HTML error</html>")
	apiClient := NewClient("test").WithBaseURL("FAKEURL").WithHTTPClient(client)
	body := NewAnalyzeSentimentsRequestBody([]string{"foo", "bar"})
	sentiments, err := apiClient.AnalyzeSentiments(Restaurants, body)
	assert.Nil(t, sentiments)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Error deserializing response")
	// Assert that this error has been wrapped
	cause := errors.Cause(err)
	assert.NotEqual(t, err, cause)

	// same for ontologies
	ontologies, err := apiClient.ListOntologies()
	assert.Nil(t, ontologies)
	assert.NotNil(t, err)
	cause = errors.Cause(err)
	assert.NotEqual(t, err, cause)
}

func TestInternalErrorStates(t *testing.T) {
	client := mocks.NewErrorClient(errors.New("Test error"))
	apiClient := NewClient("test").WithBaseURL("FAKEURL").WithHTTPClient(client)

	_, err := apiClient.get("foo")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Request failed")

	_, err = apiClient.post("bar", nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Request failed")

	// Channels cannot be serialized, and will trigger an error on json.Marshal
	bad := make(chan bool)
	defer close(bad)

	_, err = apiClient.post("baz", bad)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "JSON serialization failed")
}
