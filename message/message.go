package message

import "fmt"

type Message struct {
	Sender  string `json:"sender"`
	Payload []byte `json:"payload"`
}

func (m *Message) Log() string {
	return fmt.Sprintf("%s: %s", m.Sender, m.PayloadString())
}

func (m *Message) PayloadString() string {
	return string(m.Payload)
}