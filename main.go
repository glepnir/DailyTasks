// Package main provides ...
package main

import (
	"log"
	"net/http"

	"github.com/taigacute/DailyTasks/routers"
	"github.com/taigacute/DailyTasks/util/jsonconfig"
	"github.com/taigacute/DailyTasks/view"
)

func main() {
	values, err := jsonconfig.ReadConfig("./config/config.json")
	if err != nil {
		log.Println("Port is not allow")
	}
	view.RenderTemplate()
	routers.InitRouter()
	log.Println("Running on Port", values.ServerPort)
	log.Fatal(http.ListenAndServe(values.ServerPort, nil))
}
