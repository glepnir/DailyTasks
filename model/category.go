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

func AddCategory(username, category string) error {
	userID, err := GetUserID(username)
	if err != nil {
		return nil
	}
	log.Println("executing query to add category")
	err = database.TaskExec("insert into category(name, user_id) values(?,?)", category, userID)
	return err
}

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

//GetCategoryByName will return categoryID
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

//DeleteCategoryByName will be used to delete a category from the category page
func DeleteCategoryByName(username, category string) error {
	categoryID := GetCategoryByName(username, category)
	userID, err := GetUserID(username)
	if err != nil {
		return err
	}
	query := "update task set cat_id = null where id =? and user_id = ?"
	err = database.TaskExec(query, categoryID, userID)
	if err == nil {
		err = database.TaskExec("delete from category where id=? and user_id=?", categoryID, userID)
		if err != nil {
			return err
		}
	}
	return err
}
