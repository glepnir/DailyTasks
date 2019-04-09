// Package controllers provides ...
package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/taigacute/DailyTasks/model"
	"github.com/taigacute/DailyTasks/util/sessions"
	"github.com/taigacute/DailyTasks/view"
	"github.com/thewhitetulip/Tasks/db"
)

//AddCategoryFunc  will add category
func AddCategoryFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	r.ParseForm()
	category := r.Form.Get("category")
	if strings.Trim(category, "") != "" {
		username := sessions.GetCurrentUserName(r)
		log.Println("adding category")
		err := model.AddCategory(username, category)
		if err != nil {
			view.Message = "error adding category"
			http.Redirect(w, r, "/", http.StatusBadRequest)
		} else {
			view.Message = "Add category"
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}

//DeleteCategoryFunc will delete any category
func DeleteCategoryFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	categoryName := r.URL.Path[len("/del-category/"):]
	username := sessions.GetCurrentUserName(r)
	err := db.DeleteCategoryByName(username, categoryName)
	if err != nil {
		view.Message = "error deleting category"
	} else {
		view.Message = "Category " + categoryName + " deleted"
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
