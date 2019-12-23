package main

import (
	"github.com/websocket"
	"log"
)

type Client struct {
	socket *websocket.Conn
	send   chan []byte
	room   *Room
}

func (c *Client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			log.Println(err)
			break
		}
	}
	c.socket.Close()
}

func (c *Client) write() {
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
	c.socket.Close()
}
