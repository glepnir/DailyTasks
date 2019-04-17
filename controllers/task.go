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
	"path/filepath"
	"strings"

	"github.com/taigacute/DailyTasks/model"
	"github.com/taigacute/DailyTasks/util/redirect"
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

//AddTaskFunc Add task controller
func AddTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" { // Will work only for POST requests, will redirect to home
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	var filelink string // will store the html when we have files to be uploaded, appened to the note content
	r.ParseForm()
	file, handler, err := r.FormFile("uploadfile")
	if err != nil && handler != nil {
		//Case executed when file is uploaded and yet an error occurs
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

		if handler.Filename != "" {
			r.ParseMultipartForm(32 << 20) //defined maximum size of file
			defer file.Close()
			htmlFilename := strings.Replace(handler.Filename, " ", "-", -1)
			randomFileName := md5.New()
			io.WriteString(randomFileName, strconv.FormatInt(time.Now().Unix(), 10))
			io.WriteString(randomFileName, htmlFilename)
			token := fmt.Sprintf("%x", randomFileName.Sum(nil))
			filesDirPath := filepath.Join(".", "files")
			if err := os.MkdirAll(filesDirPath, os.ModeDir|os.ModePerm); err != nil {
				log.Println(err)
				return
			}
			f, err := os.OpenFile(filepath.Join(filesDirPath, htmlFilename), os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				log.Println(err)
				return
			}
			defer f.Close()
			io.Copy(f, file)
			if strings.HasSuffix(htmlFilename, ".png") || strings.HasSuffix(htmlFilename, ".jpg") {
				filelink = "<br> <img src='/files/" + htmlFilename + "'/>"
			} else {
				filelink = "<br> <a href=/files/" + htmlFilename + ">" + htmlFilename + "</a>"
			}
			content = content + filelink
			fileTruth := model.AddFile(htmlFilename, token, username)
			if fileTruth != nil {
				view.Message = "Error adding filename in db"
				log.Println("error adding task to db")
			}
		}
		taskTruth := tk.AddTask(title, content, category, taskPriority, username, hidden)
		if taskTruth != nil {
			view.Message = "Error adding task"
			log.Println("error adding task to db")
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		} else {
			view.Message = "Task added"
			log.Println("added task to db")
			http.Redirect(w, r, "/", http.StatusFound)
		}
	} else {
		log.Println("CSRF mismatch")
		view.Message = "Server Error"
		http.Redirect(w, r, "/", http.StatusInternalServerError)
	}
}

// UpdateTaskFunc handler
func UpdateTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	r.ParseForm()
	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		log.Println(err)
	}
	category := r.Form.Get("category")
	title := r.Form.Get("title")
	content := r.Form.Get("content")
	priority, err := strconv.Atoi(r.Form.Get("priority"))
	if err != nil {
		log.Println(err)
	}
	username := sessions.GetCurrentUserName(r)
	var hidden int
	hideTimeline := r.FormValue("hide")
	if hideTimeline != "" {
		hidden = 1
	} else {
		hidden = 0
	}
	err = tk.UpdateTask(id, title, content, category, priority, username, hidden)
	if err != nil {
		view.Message = "Error updating task"
	} else {
		view.Message = "Task updated"
		log.Println(view.Message)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

//DeleteTaskFunc is used to delete a task, trash = move to recycle bin, delete = permanent delete
func DeleteTaskFunc(w http.ResponseWriter, r *http.Request) {
	username := sessions.GetCurrentUserName(r)
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	id := r.URL.Path[len("/delete/"):]
	if id == "all" {
		err := model.DeleteAll(username)
		if err != nil {
			view.Message = "Error delete tasks"
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		id, err := strconv.Atoi(id)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/", http.StatusBadRequest)
		} else {
			err = model.DeleteTask(username, id)
			if err != nil {
				view.Message = "Error deleting task"
			} else {
				view.Message = "Task deleted"
			}
			http.Redirect(w, r, "/deleted", http.StatusFound)
		}
	}
}

//RestoreTaskFunc is used to restore task from trash, handles "/restore/" URL
func RestoreTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(r.URL.Path[len("/restore/"):])
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/deleted", http.StatusBadRequest)
	} else {
		username := sessions.GetCurrentUserName(r)
		err = tk.RestoreTask(username, id)
		if err != nil {
			view.Message = "Restore failed"
		} else {
			view.Message = "Task restored"
		}
		http.Redirect(w, r, "/deleted/", http.StatusFound)
	}
}

//RestoreFromCompleteFunc restores the task from complete to pending
func RestoreFromCompleteFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(r.URL.Path[len("/incomplete/"):])
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/completed", http.StatusBadRequest)
	} else {
		username := sessions.GetCurrentUserName(r)
		err = tk.RestoreTaskFromComplete(username, id)
		if err != nil {
			view.Message = "Restore failed"
		} else {
			view.Message = "Task restored"
		}
		http.Redirect(w, r, "/completed", http.StatusFound)
	}
}

//ShowTrashTaskFunc is used to handle the "/trash" URL which is used to show the deleted tasks
func ShowTrashTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		username := sessions.GetCurrentUserName(r)
		categories := model.GetCategories(username)
		context, err := tk.GetAllTasks(username, "deleted", "")
		context.Categories = categories
		if err != nil {
			http.Redirect(w, r, "/trash", http.StatusInternalServerError)
		}
		if view.Message != "" {
			context.Message = view.Message
			view.Message = ""
		}
		err = view.DeletedTemplate.Execute(w, context)
		if err != nil {
			log.Fatal(err)
		}
	}
}

//ShowCompleteTasksFunc is used to populate the "/completed/" URL
func ShowCompleteTasksFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		username := sessions.GetCurrentUserName(r)
		categories := model.GetCategories(username)
		context, err := tk.GetAllTasks(username, "completed", "")
		context.Categories = categories
		if err != nil {
			http.Redirect(w, r, "/completed", http.StatusInternalServerError)
		}
		view.CompletedTemplate.Execute(w, context)
	}
}

//CompleteTaskFunc is used to show the complete tasks, handles "/completed/" url
func CompleteTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	redirectURL := redirect.GetRedirectUrl(r.Referer())
	id, err := strconv.Atoi(r.URL.Path[len("/complete/"):])
	if err != nil {
		log.Println(err)
	} else {
		username := sessions.GetCurrentUserName(r)
		err = tk.CompleteTask(username, id)
		if err != nil {
			view.Message = "Complete task failed"
		} else {
			view.Message = "Task marked complete"
		}
		http.Redirect(w, r, redirectURL, http.StatusFound)
	}
}

//TrashTaskFunc is used to populate the trash tasks
func TrashTaskFunc(w http.ResponseWriter, r *http.Request) {
	//for best UX we want the user to be returned to the page making
	//the delete transaction, we use the r.Referer() function to get the link
	redirectURL := redirect.GetRedirectUrl(r.Referer())

	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.URL.Path[len("/trash/"):])
	if err != nil {
		log.Println("TrashTaskFunc", err)
		view.Message = "Incorrect command"
		http.Redirect(w, r, redirectURL, http.StatusFound)
	} else {
		username := sessions.GetCurrentUserName(r)
		err = tk.TrashTask(username, id)
		if err != nil {
			view.Message = "Error trashing task"
		} else {
			view.Message = "Task trashed"
		}
		http.Redirect(w, r, redirectURL, http.StatusFound)

	}
}

//EditTaskFunc is used to edit tasks, handles "/edit/" URL
func EditTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.URL.Path[len("/edit/"):])
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	redirectURL := redirect.GetRedirectUrl(r.Referer())
	username := sessions.GetCurrentUserName(r)
	task, err := model.GetTaskByID(username, id)
	categories := model.GetCategories(username)
	task.Categories = categories
	task.Referer = redirectURL

	if err != nil {
		task.Message = "Error fetching Tasks"
	}
	view.EditTemplate.Execute(w, task)
}

//SearchTaskFunc is used to handle the /search/ url, handles the search function
func SearchTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	r.ParseForm()
	query := r.Form.Get("query")
	username := sessions.GetCurrentUserName(r)
	context, err := model.SearchTask(username, query)
	if err != nil {
		log.Println("error fetching search results")
	}
	categories := model.GetCategories(username)
	context.Categories = categories
	view.SearchTemplate.Execute(w, context)
}
