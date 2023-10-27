package main

import (
	"github.com/dzeleniak/chatroom-gotcp/server"
)

func main() {
	s, _ := server.New()

	s.Listen(8080)
}
