// Package controllers provides ...
package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/taigacute/DailyTasks/model"
	"github.com/taigacute/DailyTasks/util/sessions"
	"github.com/taigacute/DailyTasks/view"
)

var (
	tk = model.Task{}
)

//ShowAllTasksFunc is uesd to handle the "/" URL
func ShowAllTasksFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		username := sessions.GetCurrentUserName(r)
		context, err := tk.GetAllTasks(username, "pending", "")
		log.Println(context)
		categories := model.GetCategories(username)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		} else {
			if view.Message != "" {
				context.Message = view.Message
			}
			context.CSRFToken = "abcd"
			context.Categories = categories
			view.Message = ""
			expiration := time.Now().Add(365 * 24 * time.Hour)
			cookie := http.Cookie{
				Name:    "csrftoken",
				Value:   "abcd",
				Expires: expiration,
			}
			http.SetCookie(w, &cookie)
			view.HomeTemplate.Execute(w, context)
		}
	}
}
