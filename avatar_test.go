package main

import (
	"testing"

	gomniauthtest "github.com/stretchr/gomniauth/test"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	testUser := &gomniauthtest.TestUser{}
	testUser.On("AvatarURL").Return("", ErrNoAvatarURL)
	testChatUser := &chatUser{User: testUser}
	_, err := authAvatar.GetAvatarURL(testChatUser)
	if err != ErrNoAvatarURL {
		t.Error("AuthAvatar.GetAvatarURL should return ErrNoAvatarURL" +
			"when no value present")
	}
	testUrl := "http://url-to-gravatar/"
	testUser = &gomniauthtest.TestUser{}
	testChatUser.User = testUser
	testUser.On("AvatarURL").Return(testUrl, nil)
	url, err := authAvatar.GetAvatarURL(testChatUser)
	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should return no error" +
			"when value present")
	}
	if url != testUrl {
		t.Error("AuthAvatar.GetAvatarURL should return correct URL")
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	user := &chatUser{uniqueID: "abc"}
	url, err := gravatarAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURL should not return an error")
	}
	if url != "//www.gravatar.com/avatar/abc" {
		t.Errorf("GravatarAvatar.GetAvatarURL wrongly returned %s", url)
	}
}
