package main

import "time"

// Represents a message.
type message struct {
	Name      string
	Message   string
	When      time.Time
	AvatarURL string
}
