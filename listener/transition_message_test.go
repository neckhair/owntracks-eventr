package listener_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/neckhair/owntracks-eventr/listener"
)

func TestNewTransitionMessage(t *testing.T) {
	payload := []byte("{\"Tst\":1234,\"Event\":\"enter\"}")

	tm, err := listener.NewTransitionMessage(payload)

	assert.Nil(t, err)

	assert.EqualValues(t, 1234, tm.Tst)
	assert.Equal(t, "enter", tm.Event)
}

func TestNewTransitionMessageWithEmptyPayload(t *testing.T) {
	payload := []byte("")

	tm, err := listener.NewTransitionMessage(payload)

	assert.IsType(t, &json.SyntaxError{}, err)
	assert.IsType(t, &listener.TransitionMessage{}, tm)
}

func TestTimestamp(t *testing.T) {
	tm := listener.TransitionMessage{Tst: 1489434339}
	timestamp := tm.Timestamp()

	assert.EqualValues(t, 1489434339, timestamp.Unix())
}

func TestTimestampWithUninitializedTst(t *testing.T) {
	tm := listener.TransitionMessage{}
	timestamp := tm.Timestamp()

	assert.EqualValues(t, 0, timestamp.Unix())
}
