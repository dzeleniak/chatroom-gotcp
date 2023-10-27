package user

import (
	"net"
	"testing"
)

func TestUser_SendMessageSuccess(t *testing.T) {
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

	_, err := user.SendMessage()
	if err != nil {
		t.Error(err.Error())
	}
}