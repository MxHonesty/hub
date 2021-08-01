package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// Represents a chatting user.
type client struct {
	// Websocket for the client
	socket *websocket.Conn
	// Sends messages to websocket
	send chan *message
	room *room
	// Map of cookie data.
	userData map[string]interface{}
}

// Loop for reading messages from the websocket.
// Forwards messages to room.
func (c *client) read() {
	defer c.socket.Close()
	for {
		// Reads message from websocket and forwards it to room.
		var msg *message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			return
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		c.room.forward <- msg
	}
}

// Loop for writing messages to the websocket.
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			break
		}
	}
}
