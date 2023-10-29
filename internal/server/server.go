package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/dzeleniak/chatroom-gotcp/internal/user"
	"github.com/gofrs/uuid"
)

type Server struct {
	Listener net.Listener;
	Users *sync.Map;
}	

func New() (*Server, error) {
	return &Server{
		Users: &sync.Map{},
	}, nil
}

func (s *Server) Listen(port int) (*Server, error) {
	// Start listening on the server
	networkAddress := ":" + fmt.Sprint(port);
	l, err := net.Listen("tcp",  networkAddress);

	if err != nil {
		return s, err;
	}

	s.Listener = l;

	// Listener runs on loop accepting and handling new connections
	for {
		u, err := s.acceptConnection()

		if err != nil {
			log.Printf("Error: %s\n", err)
		}

		go s.handleConnection(u);
	}
}

func (s *Server) acceptConnection() (*user.User, error) {
	// Accept new connection to the server
	c, err := s.Listener.Accept()

	if err != nil {
		return &user.User{}, err
	}

	// First message on connect should be the username
	username, err := bufio.NewReader(c).ReadBytes('\n')
	
	if err != nil {
		return &user.User{}, err;
	}
	
	log.Printf("Register Connection: %s -> %s\n", c.RemoteAddr().String(), username);

	// Generate id for the new user
	userId, err := uuid.NewV7()
	if err != nil {
		return &user.User{}, nil
	}

	// Create user and store in memory
	connectedUser := &user.User{
		Username: string(username[:len(username)-1]),
		Conn: c,
		Id: userId,
	}

	s.Users.Store(userId, connectedUser)
	return connectedUser, nil;
}

func(s *Server) handleConnection(u *user.User) {
	go func() {
		// Handle connection cleanup
		defer func() {
			u.Conn.Close();
			s.Users.Delete(u.Id);
		}()

		for {
			// TODO: This needs to call the User Receive
			buf, err := u.ReadMessage()
			if err != nil {
				s.Users.Delete(u.Id);
				log.Printf("Error: %s\n", err.Error());
				return;
			}

			
			if err != nil {
				s.Users.Delete(u.Id)
			}

			log.Println(buf);

			s.Users.Range(func(k, v interface{}) bool {
				if user, ok := v.(*user.User); ok {
					if user.Id != u.Id {
						if _, err := user.WriteMessage(buf); err != nil {
							log.Printf("Error: %s\n", err)
						}
					}
				}
				return true;
			})
		}	
	}()
}