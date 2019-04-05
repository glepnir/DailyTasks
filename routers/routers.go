// Package routers provides ...
package routers

import (
	"net/http"

	. "github.com/taigacute/DailyTasks/controllers"
	"github.com/taigacute/DailyTasks/middleware"
)

//InitRouter bind url with handler
func InitRouter() {
	http.HandleFunc("/login/", LoginFunc)
	http.HandleFunc("/loginout/", middleware.RequiresLogin(LoginOutFunc))
	http.HandleFunc("/signup/", SignUpFunc)
	http.Handle("/static/", http.FileServer(http.Dir("public")))
}
