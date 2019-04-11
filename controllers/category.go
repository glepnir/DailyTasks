// Package controllers provides ...
package controllers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/taigacute/DailyTasks/model"
	"github.com/taigacute/DailyTasks/util/sessions"
	"github.com/taigacute/DailyTasks/view"
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

func UpdateCategoryFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	var redirectURL string
	r.ParseForm()
	oldName := r.URL.Path[len("/upd-category/"):]
	newName := r.Form.Get("catname")
	username := sessions.GetCurrentUserName(r)
	err := model.UpdateCategoryByName(username, oldName, newName)
	if err != nil {
		view.Message = "error updating category"
		log.Println("not updated category " + oldName)
		redirectURL = "/category/" + oldName
	} else {
		view.Message = "cat " + oldName + " -> " + newName
		redirectURL = "/category/" + newName
	}
	log.Println("redirecting to " + redirectURL)
	http.Redirect(w, r, redirectURL, http.StatusFound)

}

//DeleteCategoryFunc will delete any category
func DeleteCategoryFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	categoryName := r.URL.Path[len("/del-category/"):]
	username := sessions.GetCurrentUserName(r)
	err := model.DeleteCategoryByName(username, categoryName)
	if err != nil {
		view.Message = "error deleting category"
	} else {
		view.Message = "Category " + categoryName + " deleted"
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

//ShowCategoryFunc will populate the /category/<id> URL which shows all the tasks related
// to that particular category
func ShowCategoryFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" && sessions.IsLoggedIn(r) {
		category := r.URL.Path[len("/category/"):]
		username := sessions.GetCurrentUserName(r)
		context, err := tk.GetAllTasks(username, "", category)
		categories := model.GetCategories(username)

		if err != nil {
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
		if view.Message != "" {
			context.Message = view.Message
		}
		context.CSRFToken = "abcd"
		context.Categories = categories
		view.Message = ""
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "csrftoken", Value: "abcd", Expires: expiration}
		http.SetCookie(w, &cookie)
		view.HomeTemplate.Execute(w, context)
	}
}
