// Package midedleware provides ...
package middleware

import (
	"net/http"

	"github.com/taigacute/DailyTasks/util/sessions"
)

//RequiresLogin is a midedleware which will be used for each
// htppHandlerFunc to check if there is any active session
func RequiresLogin(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if sessions.IsLoggedIn(r) {
			http.Redirect(w, r, "/login", 302)
			return
		}
		handler.ServeHTTP(w, r)
	}
}
