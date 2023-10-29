package user

import (
	"bufio"
	"encoding/json"
	"net"

	"github.com/dzeleniak/chatroom-gotcp/internal/message"
	"github.com/gofrs/uuid"
)

type User struct {
	Username string
	Conn     net.Conn
	Id 			uuid.UUID
}

// Read from this users connection and return a buffer containing
// the marshalled json of the message
func (u *User) ReadMessage() ([]byte, error) {

	buf, err := bufio.NewReader(u.Conn).ReadBytes('\n')
	if err != nil {
		return nil, err;
	}

	message := &message.Message{
		Sender: u.Username,
		Payload: buf,
	}

	bytes, err := json.Marshal(message)

	if err != nil {
		return nil, err;
	}

	return bytes, nil;
}

// Receives message from the server to the users connection
func (u *User) WriteMessage(buf []byte) (*message.Message, error) {
	var msg *message.Message;
	if err := json.Unmarshal(buf, &msg); err != nil {
		return &message.Message{}, err;
	}

	if _, err := u.Conn.Write([]byte(msg.Log())); err != nil {
		return &message.Message{}, err;
	}

	return msg, nil;
}