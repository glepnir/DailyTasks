// Package database provides the method of database option
// this package Encapsulation database/sql
package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var database Database
var err error

//Database encapsulates database
type Database struct {
	db *sql.DB
}

//Begin() will provide tx
func (d Database) begin() (tx *sql.Tx) {
	tx, err := d.db.Begin()
	if err != nil {
		log.Println(err)
		return nil
	}
	return tx
}

// Prepare() Encapsulation datbase/sql
func (d Database) prepare(q string) (stmt *sql.Stmt) {
	stmt, err := d.db.Prepare(q)
	if err != nil {
		log.Println(err)
		return nil
	}
	return stmt
}

//Query() Encapsulation database/sql
func (d Database) query(q string, args ...interface{}) (rows *sql.Rows) {
	rows, err := d.db.Query(q, args...)
	if err != nil {
		log.Println(err)
		return nil
	}
	return rows
}

func (d Database) queryrow(q string, args ...interface{}) (row *sql.Row) {
	return d.db.QueryRow(q, args...)
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

//TaskExec
func TaskExec(q string, args ...interface{}) error {
	SQL := database.prepare(q)
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

func TaskQueryRows(q string, args ...interface{}) (rows *sql.Rows) {
	return database.query(q, args...)
}

func TaskQueryRow(q string, args ...interface{}) (row *sql.Row) {
	return database.queryrow(q, args...)
}
