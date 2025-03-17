package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
	Rooms map[string]*websocket.Conn
}

func NewServer() *Server {
	return &Server{
		Rooms: make(map[string]*websocket.Conn),
	}
}

func (s *Server) HandleWS(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade websocket: %v", err)
		return
	}
	defer conn.Close()

	fmt.Printf("New incoming connections from: %s", conn.RemoteAddr())

	s.conns[conn] = true

	s.readLoop(conn)
}

func (s *Server) readLoop(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading message connection: %s", err)
			continue
		}

		fmt.Printf("Received message: %s", msg)
	}
}

func (s *Server) JoinRoom(room string, conn *websocket.Conn) {
	s.Rooms[room] = conn
}

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func handleJoin(conn *websocket.Conn, message Message) {
	log.Println(conn)
	log.Printf("User %s joined", message.Data)
}

func main() {
	server := NewServer()

	http.Handle("/ws", http.HandlerFunc(server.HandleWS))

	log.Println("Server started!")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
