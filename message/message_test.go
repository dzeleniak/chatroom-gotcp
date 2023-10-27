package message

import (
	"testing"
)

func TestMessage_LogCorrectOutput(t *testing.T) {
	message := &Message{
		Sender: "dzeleniak",
		Payload: []byte("Hello from me!"),
	}

	res := message.Log();
	expectation := "dzeleniak: Hello from me!"

	if res != expectation {
		t.Errorf("Expected: %s\tReceived: %s", expectation, res)
	}
}