// Package sessions provides ...
package sessions

import (
	"net/http"

	"log"

	"github.com/gorilla/sessions"
)

// Store the cookie Store
var Store = sessions.NewCookieStore([]byte("tasks-g0lang"))

//IsLoggedIn will check if the user has an active sesion and return true
func IsLoggedIn(r *http.Request) bool {
	session, err := Store.Get(r, "session")
	if err != nil {
		log.Println(err)
	}
	if session.Values["loggedin"] == "true" {
		return true
	}
	return false
}

//GetCurrentUserName will return current username
func GetCurrentUserName(r *http.Request) string {
	session, err := Store.Get(r, "session")
	if err == nil {
		return session.Values["username"].(string)
	}
	return ""
}
