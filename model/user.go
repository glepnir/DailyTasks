// Package model provides ...
package model

import (
	"github.com/taigacute/DailyTasks/database"
)

//User struct
type User struct {
	id       int
	Username string
	Password string
}

//RegisterUser add  user
func (user User) RegisterUser(uname string, pwd string, email string) error {
	sql := "insert into user(username,password,email)values(?,?,?)"
	err := database.TaskQuery(sql, uname, pwd, email)
	return err
}
