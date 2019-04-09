// Package model provides ...
package model

import (
	"time"

	"log"

	"github.com/taigacute/DailyTasks/database"
)

//Comment is the struct used to populate comments per tasks
type Comment struct {
	ID       int    `json:"id"`
	Content  string `json:"content"`
	Created  string `json:"created_date"`
	Username string `json:"username"`
}

//AddComment  will add comment
func AddComment(username string, id int, comment string) error {
	userID, err := GetUserID(username)
	if err != nil {
		return err
	}
	stmt := "insert into comments(taskID, content, created, user_id) values (?,?,datetime(),?)"
	err = database.TaskExec(stmt, id, comment, userID)
	if err != nil {
		return err
	}
	log.Println("add comment to task ID", id)
	return nil
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

//DeleteCommentByID will actually delete the comment from db
func DeleteCommentByID(username string, id int) error {
	userID, err := GetUserID(username)
	if err != nil {
		return err
	}
	query := "delete from comments where id=? and user_id=?"
	err = database.TaskExec(query, id, userID)
	return err
}
