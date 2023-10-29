package main

import (
	"github.com/dzeleniak/chatroom-gotcp/pkg/server"
)

func main() {
	s, _ := server.New()
	s.Listen(8080)
}
