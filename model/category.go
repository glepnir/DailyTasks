// Package model provides...
package model

import (
	"log"

	"github.com/taigacute/DailyTasks/database"
)

//CategoryCount is the struct used to populate the sidebar
//which contains the category name and the count of the tasks
//in each category
type CategoryCount struct {
	Name  string
	Count int
}

//Category is the structure of the category table
type Category struct {
	ID      int    `json:"category_id"`
	Name    string `json:"category_name"`
	Created string `json:"created_date"`
}

//Categories will show
type Categories []Category

//GetCategories will return the list of cateories
//render in the template
func GetCategories(uname string) []CategoryCount {
	userID, err := GetUserID(uname)
	if err != nil {
		return nil
	}
	stmt := "select 'UNCATEGORIZED' as name, count(1) from task where cat_id=0 union  select c.name, count(*) from   category c left outer join task t  join status s on  c.id = t.cat_id and t.task_status_id=s.id where s.status!='DELETED' and c.user_id=?   group by name    union     select name, 0  from category c, user u where c.user_id=? and name not in (select distinct name from task t join category c join status s on s.id = t.task_status_id and t.cat_id = c.id and s.status!='DELETED' and c.user_id=?)"
	rows := database.TaskQueryRows(stmt, userID, userID, userID)
	var cateories []CategoryCount
	var category CategoryCount
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&category.Name, &category.Count)
		if err != nil {
			log.Println(err)
		}
		cateories = append(cateories, category)
	}
	return cateories
}

func GetCategoryByName(username, category string) int {
	stmt := "select id from category where name=? and user_id = (select id from user where username=?)"
	rows := database.TaskQueryRows(stmt, category, username)
	var categoryID int
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&categoryID)
		if err != nil {
			log.Println(err)
		}
	}
	return categoryID
}
