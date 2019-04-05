// Package model provides ...
package model

//Context is the struct passed to templates
type Context struct {
	Tasks      []Task
	Navigation string
	Search     string
	Message    string
	CSRFToken  string
	Categories []CategoryCount
	Referer    string
}
