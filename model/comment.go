// Package model provides ...
package model

import (
	"time"

	"github.com/taigacute/DailyTasks/database"
)

//Comment is the struct used to populate comments per tasks
type Comment struct {
	ID       int    `json:"id"`
	Content  string `json:"content"`
	Created  string `json:"created_date"`
	Username string `json:"username"`
}

//GetComments will return a map and error which key is int value is slice of comment
func GetComments(username string) (map[int][]Comment, error) {
	commentMap := make(map[int][]Comment)
	var taskID int
	var comment Comment
	var created time.Time
	userID, err := GetUserID(username)
	if err != nil {
		return commentMap, err
	}
	stmt := "select c.id,c.taskID,c.content,c.created,u.username from comments c ,task t, user u where t.id=c.taskID and c.user_id=t.user_id and t.user_id=u.id and u.id=?"
	rows := database.TaskQueryRows(stmt, userID)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&comment.ID, &taskID, &comment.Content, &created, &comment.Username)
		if err != nil {
			return commentMap, err
		}
		created = created.Local()
		comment.Created = created.Format("Jan 2 2006 15:04:05")
		commentMap[taskID] = append(commentMap[taskID], comment)
	}
	return commentMap, nil
}
