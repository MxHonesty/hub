package main

import (
	"fmt"
	"hub/trace"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

// Represents a chat room.
type room struct {
	// Holds incoming messages
	forward chan *message
	join    chan *client
	leave   chan *client
	clients map[*client]bool // All current clients in the room
	tracer  trace.Tracer     // Trace room activity
}

// Returns a new room.
// The default tracer is turned Off.
func newRoom() *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}

// Main room loop.
func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// joining
			r.clients[client] = true
			r.tracer.Trace("[+] Client joined room")
		case client := <-r.leave:
			// leaving
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("[-] Client left room")
		case msg := <-r.forward:
			// forward message to all clients
			for client := range r.clients {
				client.send <- msg
				r.tracer.Trace(fmt.Sprintf("[M] %s sent '%s' at %s.",
					msg.Name, msg.Message, msg.When.String()))
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize,
	WriteBufferSize: socketBufferSize}

// Init websocket connection.
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP: ", err)
		return
	}
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie: ", err)
		return
	}
	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
