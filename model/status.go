// Package model provides ...
package model

//Status is the JSON struct to be returned
type Status struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}
