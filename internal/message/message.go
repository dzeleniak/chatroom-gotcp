package message

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	Sender  string `json:"sender"`
	Payload []byte `json:"payload"`
}

func New(username string, payload []byte) *Message {
	
	return &Message{
		Sender: username,
		Payload: payload,
	}

}

func (m *Message) Log() string {
	return fmt.Sprintf("%s: %s", m.Sender, m.PayloadString())
}

func (m *Message) PayloadString() string {
	return string(m.Payload)
}

func ToBytes(msg *Message) ([]byte, error){
	return json.Marshal(msg)
}

func FromBytes(msg []byte) (*Message, error) {
	var resultMessage *Message;
	err := json.Unmarshal(msg, resultMessage)

	return resultMessage, err;
}