// Package routers provides ...
package routers

import (
	"net/http"

	. "github.com/taigacute/DailyTasks/controllers"
	"github.com/taigacute/DailyTasks/middleware"
)

//InitRouter bind url with handler
func InitRouter() {
	//login logout signup
	http.HandleFunc("/login/", LoginFunc)
	http.HandleFunc("/logout/", middleware.RequiresLogin(LoginOutFunc))
	http.HandleFunc("/signup/", SignUpFunc)
	http.Handle("/static/", http.FileServer(http.Dir("public")))

	//these handlers fetch set of tasks
	http.HandleFunc("/", middleware.RequiresLogin(ShowAllTasksFunc))
	http.HandleFunc("/add/", middleware.RequiresLogin(AddTaskFunc))
	http.HandleFunc("/add-category/", middleware.RequiresLogin(AddCategoryFunc))
	http.HandleFunc("/add-comment/", middleware.RequiresLogin(AddCommentFunc))

	http.HandleFunc("/del-comment/", middleware.RequiresLogin(DeleteCommentFunc))
	http.HandleFunc("/del-category/", middleware.RequiresLogin(DeleteCategoryFunc))
	http.HandleFunc("/delete/", middleware.RequiresLogin(DeleteTaskFunc))
}
