package main

import (
	"errors"
)

// Error that is returned when Avatar instance is unable to
// provide an Avatar URL.
var ErrNoAvatarURL = errors.New("hub: Unable to get an avatar URL")

// Avatar common interface for user profile picture
// representations.
type Avatar interface {
	// Gets avatar URL for the specified client.
	// Retunrs an error if that cannot be done.
	//	Errors:
	//		ErrNoAvatarURL is returned if the object is unable
	// 		to get a URL for the specified client.
	GetAvatarURL(ChatUser) (string, error)
}

// Chain of Responsabiltiy for determining the used Avatar type
type TryAvatars []Avatar

func (a TryAvatars) GetAvatarURL(c ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(c); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(c ChatUser) (string, error) {
	url := c.AvatarURL()
	if len(url) == 0 {
		return "", ErrNoAvatarURL
	}
	return url, nil
}

type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(c ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + c.UniqueID(), nil
}
