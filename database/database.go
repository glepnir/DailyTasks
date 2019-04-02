// Package database provides the method of database option
// this package Encapsulation database/sql
package database

import (
	"database/sql"
	"log"
)

var database Database
var err error

//Database encapsulates database
type Database struct {
	db *sql.DB
}

//Encapsulation the database/sql begin()
func (d Database) begin() (tx *sql.Tx) {
	tx, err := d.db.Begin()
	if err != nil {
		log.Println(err)
		return nil
	}
	return tx
}

//Encapsulation datbase/sql Prepare()
func (d Database) prepare(q string) (stmt *sql.Stmt) {
	stmt, err := d.db.Prepare(q)
	if err != nil {
		log.Println(err)
		return nil
	}
	return stmt
}

//Encapsulation database/sql Query()
func (d Database) query(q string, args ...interface{}) (rows *sql.Rows) {
	rows, err := d.db.Query(q, args...)
	if err != nil {
		log.Println(err)
		return nil
	}
	return rows
}

func init() {
	database.db, err = sql.Open("sqlite3", "./config/tasks.db")
	if err != nil {
		log.Fatal(err)
	}
}

//Close function closes this database connection
func Close() {
	database.db.Close()
}

func TaskQuery(sql string, args ...interface{}) error {
	SQL := database.prepare(sql)
	tx := database.begin()
	_, err = tx.Stmt(SQL).Exec(args...)
	if err != nil {
		log.Println("TaskQuery: ", err)
		tx.Rollback()
	} else {
		err = tx.Commit()
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("Commit successful")
	}
	return err
}
