// Package controllers provides ...
package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"text/template"

	"crypto/md5"
	"fmt"
	"io"
	"os"
	"strings"

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

func AddTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	var filelink string
	r.ParseForm()
	file, handler, err := r.FormFile("uploadfile")
	if err != nil && handler != nil {
		log.Println(err)
		view.Message = "Error uploading file"
		http.Redirect(w, r, "/", http.StatusInternalServerError)
	}
	taskPriority, priorityErr := strconv.Atoi(r.FormValue("priority"))
	if priorityErr != nil {
		log.Print(priorityErr)
		view.Message = "Bad task priority"
	}
	priorityList := []int{1, 2, 3}
	found := false
	for _, priority := range priorityList {
		if taskPriority == priority {
			found = true
		}
	}
	if !found {
		taskPriority = 1
	}
	var hidden int
	hideTimeline := r.FormValue("hide")
	if hideTimeline != "" {
		hidden = 1
	} else {
		hidden = 0
	}
	category := r.FormValue("category")
	title := template.HTMLEscapeString(r.Form.Get("title"))
	content := template.HTMLEscapeString(r.Form.Get("content"))
	formToken := template.HTMLEscapeString(r.Form.Get("CSRFToken"))
	cookie, _ := r.Cookie("csrftoken")
	if formToken == cookie.Value {
		username := sessions.GetCurrentUserName(r)
		if handler != nil {
			r.ParseMultipartForm(32 << 20)
			defer file.Close()
			htmlFilename := strings.Replace(handler.Filename, "", "-", -1)
			randomFileName := md5.New()
			io.WriteString(randomFileName, strconv.FormatInt(time.Now().Unix(), 10))
			io.WriteString(randomFileName, htmlFilename)
			token := fmt.Sprintf("%x", randomFileName.Sum(nil))
			f, err := os.OpenFile("./files/"+htmlFilename, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				log.Println(err)
				return
			}
			defer f.Close()
			if strings.HasSuffix(htmlFilename, ".png") || strings.HasSuffix(htmlFilename, ".jgp") {
				filelink = "<br> <img src='/files/" + htmlFilename + "'/>"
			} else {
				filelink = "<br> <a href=/files/" + htmlFilename + ">" + htmlFilename + "</a>"
			}
			content = content + filelink
			fileTruth := model.AddFile(htmlFilename, token, username)
			if fileTruth != nil {
				view.Message = "Error add filename in db"
				log.Println("Error add filename in db")
			}
		}
		taskTruth := tk.AddTask(title, content, category, taskPriority, username, hidden)
		if taskTruth != nil {
			view.Message = "Error adding task"
			log.Println("error adding task")
			http.Redirect(w, r, "/", http.StatusFound)
		}
	} else {
		log.Println("CSRF mismatch")
		view.Message = "Server Error"
		http.Redirect(w, r, "/", http.StatusInternalServerError)
	}
}
