// Package controllers provides ...
package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/taigacute/DailyTasks/model"
	"github.com/taigacute/DailyTasks/util/sessions"
	"github.com/taigacute/DailyTasks/view"
)

// AddCommentFunc will add a comment into db
func AddCommentFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	r.ParseForm()
	text := r.Form.Get("commentText")
	id := r.Form.Get("taskID")
	idInt, err := strconv.Atoi(id)
	if (err != nil) || (text == "") {
		log.Println("unable to convert into integer")
		view.Message = "Error adding comment"
	} else {
		username := sessions.GetCurrentUserName(r)
		err := model.AddComment(username, idInt, text)
		if err != nil {
			log.Println("unable insert into db")
			view.Message = "Comment add failed"
		} else {
			view.Message = "Comment add successful"
		}
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

//DeleteCommentFunc will delete any category
func DeleteCommentFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	id := r.URL.Path[len("/del-comment/"):]
	commentID, err := strconv.Atoi(id)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	username := sessions.GetCurrentUserName(r)
	err = model.DeleteCommentByID(username, commentID)
	if err != nil {
		view.Message = "comment not deleted"
	} else {
		view.Message = "comment deleted"
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
