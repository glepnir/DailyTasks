// Package controllers provides ...
package controllers

import (
	"log"
	"net/http"
)

// UploadedFileHandler is used to handle the uploaded file related requests
func UploadedFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	token := r.URL.Path[len("/files/"):]
	log.Println("serving file ./files/" + token)
	http.ServeFile(w, r, "./files/"+token)
}
