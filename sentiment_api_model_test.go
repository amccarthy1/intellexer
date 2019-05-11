package intellexer

import (
	"encoding/json"
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
	requestBody := AnalyzeSentimentsRequest{{ID: allZeros, Text: "foo"}}
	bytes, err := json.Marshal(requestBody)

	assert.Nil(t, err, "Request body should marshal without error")
	assert.Equal(
		t,
		`[{"id":"00000000-0000-0000-0000-000000000000","text":"foo"}]`,
		string(bytes),
	)
}
