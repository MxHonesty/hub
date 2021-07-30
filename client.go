package main

import (
	"github.com/gorilla/websocket"
)

// Represents a chatting user.
type client struct {
	socket *websocket.Conn
	send   chan []byte // sends messages to websocket
	room   *room
}

// Loop for reading messages from the websocket.
// Forwards messages to room.
func (c *client) read() {
	defer c.socket.Close()
	for {
		// Reads message from websocket and forwards it to room.
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

// Loop for writing messages to the websocket.
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
