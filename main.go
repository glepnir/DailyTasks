// Package main provides ...
package main

import (
	"log"
	"net/http"

	"github.com/taigacute/DailyTasks/routers"
	"github.com/taigacute/DailyTasks/view"
)

func main() {
	view.RenderTemplate()
	routers.InitRouter()
	log.Print("Running on Port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
