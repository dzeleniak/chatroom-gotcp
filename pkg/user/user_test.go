package user

import (
	"bufio"
	"encoding/json"
	"net"
	"testing"

	m "github.com/dzeleniak/chatroom-gotcp/pkg/message"
)

func TestUser_ReadMessageSuccess(t *testing.T) {
	server, client := net.Pipe();
	defer client.Close();
	defer server.Close();

	user := &User {
		Conn: client,
		Username: "testing",
	}

	message := "Test message :))\n"
	
	go func () {
		server.Write([]byte(message))
	}()

	buf, err := user.ReadMessage()
	if err != nil {
		t.Error(err.Error())
	}

	var resMsg *m.Message;
	err = json.Unmarshal(buf, &resMsg);

	if err != nil {
		t.Error(err.Error())
	}

	if string(resMsg.Payload) != message {
		t.Errorf("Error: Result payload not equal to original message...\nOriginal: %s\nReceived: %s\n", 
			message,
			string(resMsg.Payload))
	}

	if resMsg.Sender != user.Username {
		t.Errorf("Error: Incorrect username transmitted with message\n")
	}
}

func TestUser_WriteMessageSuccess(t *testing.T) {
	server, client := net.Pipe();
	defer server.Close()
	defer client.Close();

	user := &User{
		Conn: client,
		Username: "testing",
	}	
	
	message := &m.Message {
		Sender: user.Username,
		Payload: []byte("Test Message :))\n"),
	}
	
	buf, err := json.Marshal(message)

	if err != nil {
		t.Errorf("Error: Failed to convert message to bytes\n")
	}
	
	go func() {
		msg, err := bufio.NewReader(server).ReadBytes('\n')
	
		if err != nil {
			t.Errorf(
				"Error: Failed to read message from user on the server.\n%s\n",
				err.Error())
		}

		t.Log(msg)		
	}()

	_, err = user.WriteMessage(buf)
	if err != nil {
		t.Errorf(
			"Error: Failed to write message from user to the server\n%s\n",
			err.Error());
	}

}