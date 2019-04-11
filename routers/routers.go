// Package routers provides ...
package routers

import (
	"net/http"

	. "github.com/taigacute/DailyTasks/controllers"
	"github.com/taigacute/DailyTasks/middleware"
)

//InitRouter bind url with handler
func InitRouter() {
	http.HandleFunc("/", middleware.RequiresLogin(ShowAllTasksFunc))
	http.HandleFunc("/login/", LoginFunc)
	http.HandleFunc("/logout/", middleware.RequiresLogin(LoginOutFunc))
	http.HandleFunc("/signup/", SignUpFunc)
	http.HandleFunc("/add/", middleware.RequiresLogin(AddTaskFunc))
	http.HandleFunc("/add-category/", middleware.RequiresLogin(AddCategoryFunc))
	http.HandleFunc("/add-comment/", middleware.RequiresLogin(AddCommentFunc))
	http.HandleFunc("/del-comment/", middleware.RequiresLogin(DeleteCommentFunc))
	http.HandleFunc("/del-category/", middleware.RequiresLogin(DeleteCategoryFunc))
	http.HandleFunc("/delete/", middleware.RequiresLogin(DeleteTaskFunc))
	http.HandleFunc("/upd-category/", middleware.RequiresLogin(UpdateCategoryFunc))
	http.HandleFunc("/update/", middleware.RequiresLogin(UpdateTaskFunc))
	http.HandleFunc("/incomplete/", middleware.RequiresLogin(RestoreFromCompleteFunc))
	http.HandleFunc("/restore/", middleware.RequiresLogin(RestoreTaskFunc))
	http.HandleFunc("/category/", middleware.RequiresLogin(ShowCategoryFunc))
	http.HandleFunc("/deleted/", middleware.RequiresLogin(ShowTrashTaskFunc))
	http.HandleFunc("/completed/", middleware.RequiresLogin(ShowCompleteTasksFunc))
	http.HandleFunc("/complete/", middleware.RequiresLogin(CompleteTaskFunc))
	http.HandleFunc("/files/", middleware.RequiresLogin(UploadedFileHandler))
	http.HandleFunc("/trash/", middleware.RequiresLogin(TrashTaskFunc))
	http.HandleFunc("/edit/", middleware.RequiresLogin(EditTaskFunc))
	http.HandleFunc("/search/", middleware.RequiresLogin(SearchTaskFunc))
	http.Handle("/static/", http.FileServer(http.Dir("public")))
}
