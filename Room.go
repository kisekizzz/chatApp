package main

import (
	"github.com/websocket"
	"log"
	"net/http"
)

const (
	SOCKET_BUFFER_SIZE  = 2048
	MESSAGE_BUFFER_SIZE = 1024
)

type Room struct {
	forward       chan []byte
	clients_join  chan *Client
	clients_leave chan *Client
	clients       map[*Client]bool
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.clients_join:
			r.clients[client] = true
		case client := <-r.clients_leave:
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			for client := range r.clients {
				select {
				case client.send <- msg:
				default:
					log.Println("errors encounted in room!")
					close(client.send)
					break
				}
			}
		}
	}
}

var upgrader = &websocket.Upgrader{ReadBufferSize: SOCKET_BUFFER_SIZE, WriteBufferSize: SOCKET_BUFFER_SIZE, CheckOrigin: func(*http.Request) bool { return true }}

func (r *Room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("room start1")
	socket, err := upgrader.Upgrade(w, req, nil)
	log.Println("room start2")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("room start3")
	client := &Client{
		socket: socket,
		room:   r,
		send:   make(chan []byte, MESSAGE_BUFFER_SIZE),
	}
	log.Println("room start4")
	r.clients_join <- client
	defer func() { r.clients_leave <- client }()
	log.Println("room start3")
	go client.write()
	log.Println("room start3")
	client.read()
	log.Println("room start3")
}

func NewRoom() *Room {
	return &Room{
		forward:       make(chan []byte),
		clients_join:  make(chan *Client),
		clients_leave: make(chan *Client),
		clients:       make(map[*Client]bool),
	}
}
