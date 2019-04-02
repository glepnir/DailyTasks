// Package main provides ...
package main

import (
	"log"
	"net/http"

	"github.com/taigacute/DailyTasks/routers"
)

func main() {
	routers.InitRouter()
	log.Print("Running on Port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
