package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
)

// Handler for authentification.
type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth")
	if err == http.ErrNoCookie || cookie.Value == "" {
		// Not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		// Other error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.next.ServeHTTP(w, r)
}

// Returns a http.Handler decorated in authHandler.
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// loginHandler handles the third-party logins.
// /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	if len(segs) < 4 {
		http.Error(w, "Error when trying to auth", http.StatusBadRequest)
		return
	}

	action := segs[2]
	provider := segs[3]

	switch action {
	case "login":
		providerO, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to get provider "+
				"%s: %s", provider, err), http.StatusBadRequest)
			return
		}
		loginUrl, err := providerO.GetBeginAuthURL(nil, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to GetBeginAuthURL "+
				"for %s: %s", provider, err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)

	case "callback":
		providerO, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to get provider "+
				"%s: %s", provider, err), http.StatusBadRequest)
			return
		}
		creds, err := providerO.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to complete auth for "+
				"%s: %s", provider, err), http.StatusInternalServerError)
			return
		}
		user, err := providerO.GetUser(creds)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to get user from "+
				"%s: %s", provider, err), http.StatusInternalServerError)
			return
		}

		authCookieValue := objx.New(map[string]interface{}{
			"name":       user.Name(),
			"avatar_url": user.AvatarURL(),
			"email":      user.Email(),
		}).MustBase64()
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/",
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}
