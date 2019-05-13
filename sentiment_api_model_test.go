package intellexer

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewSentimentsRequest(t *testing.T) {
	text1 := uuid.New().String()
	text2 := uuid.New().String()
	text3 := uuid.New().String()
	requestBody := NewAnalyzeSentimentsRequestBody([]string{text1, text2, text3})

	assert.Len(t, requestBody, 3)
	assert.Equal(t, requestBody[0].Text, text1)
	assert.Equal(t, requestBody[1].Text, text2)
	assert.Equal(t, requestBody[2].Text, text3)
}

func TestJSONSerialization(t *testing.T) {
	allZeros := uuid.UUID([16]byte{})
	requestBody := []Review{{ID: allZeros, Text: "foo"}}
	bytes, err := json.Marshal(requestBody)

	assert.Nil(t, err, "Request body should marshal without error")
	assert.Equal(
		t,
		`[{"id":"00000000-0000-0000-0000-000000000000","text":"foo"}]`,
		string(bytes),
	)
}

func TestResponseDeserialization(t *testing.T) {
	bytes, err := ioutil.ReadFile("testdata/analyze_sentiments_response.json")
	assert.Nil(t, err, "testdata file should read without error")
	var res SentimentResponse
	err = json.Unmarshal(bytes, &res)
	assert.Nil(t, err)

	assert.Len(t, res.Sentiments, 1)

	sentiment := res.Sentiments[0]
	assert.Equal(t, "3fce35a7-b41c-4b75-b564-ec438cc30755", sentiment.ID)
	assert.Equal(t, 1.7317706343552688, sentiment.SentimentWeight)
	assert.Equal(t, 1, res.SentimentsCount)
	assert.Equal(t, Restaurants, res.Ontology)

	topLevelOpinion := res.Opinions
	assert.Nil(t, topLevelOpinion.Text)
	level2Opinions := topLevelOpinion.Children
	assert.Len(t, level2Opinions, 2)

	child1 := level2Opinions[0]
	assert.Equal(t, "Drinks", *child1.Text)
	assert.Equal(t, 1, child1.F)
	assert.Equal(t, 0.0, child1.SentimentWeight)

	child2 := level2Opinions[1]
	assert.Equal(t, "Other", *child2.Text)
	assert.Equal(t, 2, child2.F)
	assert.Equal(t, 0.0, child2.SentimentWeight)

	// Skip through to leaf node
	leafOpinion := child1.Children[0].Children[0]
	assert.Equal(t, "love", *leafOpinion.Text)
	assert.Equal(t, 1, leafOpinion.F)
	assert.Equal(t, 2.8, leafOpinion.SentimentWeight)
}
