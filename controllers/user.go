// Package controller provides ...
package controller

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/taigacute/DailyTasks/model"
	"github.com/taigacute/DailyTasks/view"
)

var user = model.User{}

//LoginFunc implements the login functionality
func LoginFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		view.LoginTemplate.Execute(w, nil)
	case "POST":
		r.ParseForm()
		username := r.Form.Get("username")
		password := r.Form.Get("password")
		// there will not handle the empty value it should be handle by javascript
		if user.UserIsExist(username) {
			if user.ValidUser(username, password) {
				http.Redirect(w, r, "/", 302)
			}
			http.Error(w, "Wrong username or password", http.StatusInternalServerError)
		} else {
			http.Error(w, "User doesnt exist", http.StatusInternalServerError)
		}
		view.LoginTemplate.Execute(w, nil)
	}
}

// SignUpFunc will enable new users to sign up
func SignUpFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	email := r.Form.Get("email")
	if user.UserIsExist(username) {
		http.Error(w, "UserName has exist", http.StatusInternalServerError)
	} else {
		hashpwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "unable to Generatepasswrod", http.StatusInternalServerError)
		}
		pwd := string(hashpwd)
		err = user.RegisterUser(username, pwd, email)
		if err != nil {
			http.Error(w, "Unable to sign user up", http.StatusInternalServerError)
		} else {
			http.Redirect(w, r, "/login", 302)
		}
	}
}
