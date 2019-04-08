// Package view provides ...
package view

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	HomeTemplate      *template.Template
	DeletedTemplate   *template.Template
	CompletedTemplate *template.Template
	EditTemplate      *template.Template
	SearchTemplate    *template.Template
	templates         *template.Template
	LoginTemplate     *template.Template
	Message           string
	err               error
)

//RenderTemplate  will populate tempalte
func RenderTemplate(tpl string) {
	var templateFiles []string
	var templateDir string
	switch tpl {
	case "English":
		templateDir = "./templates/"
	case "Chinese":
		templateDir = "./templatescn/"
	}
	files, err := ioutil.ReadDir(templateDir)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".html") {
			templateFiles = append(templateFiles, templateDir+filename)
		}
	}
	templates, err = template.ParseFiles(templateFiles...)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	HomeTemplate = templates.Lookup("home.html")
	DeletedTemplate = templates.Lookup("deleted.html")
	EditTemplate = templates.Lookup("edit.html")
	SearchTemplate = templates.Lookup("search.html")
	CompletedTemplate = templates.Lookup("completed.html")
	LoginTemplate = templates.Lookup("login.html")
}
