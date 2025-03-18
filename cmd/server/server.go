package main

import (
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type Server struct {
	conns map[*websocket.Conn]bool
	mu sync.Mutex
	//Rooms map[string][]*websocket.Conn
}

func NewServer() *Server {
	return &Server{
		//Rooms: make(map[string][]*websocket.Conn),
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) HandleWS(conn *websocket.Conn) {
	log.Printf("New incoming connection from: %s", conn.RemoteAddr())
	s.mu.Lock()
	s.conns[conn] = true
	s.mu.Unlock()

	defer s.removeConnection(conn)

	for {
		var msg Message
		err := websocket.JSON.Receive(conn, &msg)
		if err != nil {
			log.Printf("Error reading message connection: %s", err)
			continue
		}

		log.Printf("Received WebRTC Signal: %s from %s", msg.Type, conn.RemoteAddr())
		s.broadcast(msg, conn)
	}
}

func (s *Server) broadcast(message Message, sender *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for conn, _ := range s.conns {
		if conn != sender {
			go func(c *websocket.Conn) {
				err := websocket.JSON.Send(c, message)
				if err != nil {
					log.Printf("Error broadcasting message: %v", err)
					s.removeConnection(c)
				}
			}(conn)
		}
	}
}

func (s *Server) removeConnection(conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.conns, conn)
	conn.Close()
	log.Printf("Connection closed: %s\n", conn.RemoteAddr())
}

//func (s *Server) joinRoom(room string, conn *websocket.Conn) {
//s.mu.Lock()
//defer s.mu.Unlock()
//
//if _, ok := s.Rooms[room]; !ok {
//s.Rooms[room] = append(s.Rooms[room], conn)
//}
//
//if _, ok := s.conns[conn]; !ok {
//s.conns[conn] = room
//}
//}

func main() {
	server := NewServer()

	http.Handle("/ws", websocket.Handler(server.HandleWS))

	log.Println("Server started!")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
