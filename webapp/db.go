package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

const (
	username = "root"
	password = "pass1234"
	hostname = "0.0.0.0:3306"
	dbname   = "webappdb"
)

func dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)
}

func openDB() {
	var err error
	db, err = sql.Open("mysql", dsn())
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		panic(err)
	}
}

func closeDB() {
	db.Close()
}

func createDb() {
	openDB()

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return
	}

	no, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return
	}
	log.Printf("rows affected %d\n", no)

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)

	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return
	}
	log.Printf("Connected to DB %s successfully\n", dbname)
}

func createTable() {
	_, err := db.Exec("DROP TABLE IF EXISTS webappdb.users")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE webappdb.users(
		id varchar(100) NOT NULL,
		firstname varchar(100) COLLATE utf8_unicode_ci NOT NULL,
		lastname varchar(100) COLLATE utf8_unicode_ci NOT NULL,
		username varchar(100) COLLATE utf8_unicode_ci NOT NULL,
		password varchar(255) COLLATE utf8_unicode_ci NOT NULL,
		created datetime NOT NULL,
		modified datetime NOT NULL,
		PRIMARY KEY (id))ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;`)
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT * FROM webappdb.users")
	if err != nil {
		panic(err)
	}

	for rows.Next() {

		// Get column names
		columns, err := rows.Columns()
		if err != nil {
			panic(err.Error())
		}

		// Create interface set
		values := make([]interface{}, len(columns))
		scanArgs := make([]interface{}, len(values))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		// Scan for arbitrary values
		err = rows.Scan(scanArgs...)
		if err == nil {

			// Print data
			for i, value := range values {
				switch value.(type) {
				default:
					fmt.Printf("%s :: %s :: %+v\n", columns[i], reflect.TypeOf(value), value)
				}
			}
		} else {
			panic(err)
		}
	}
}

func queryById(id string) *User {
	user := User{}
	err := db.QueryRow(`SELECT id, firstname, lastname, username, created, modified 
							FROM webappdb.users WHERE id = ?`, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username,
		&user.AccountCreated, &user.AccountUpdated)
	if err != nil {
		log.Printf(err.Error())
		return nil
	}

	return &user
}

func insertUser(user User) bool {
	insert, err := db.Prepare(`INSERT INTO webappdb.users(id, firstname, lastname, username, password, created, modified) 
						VALUES (?, ?, ?, ?, ?, ?, ?)`)

	if err != nil {
		log.Printf(err.Error())
		return false
	}

	_, err = insert.Exec(user.ID, user.FirstName, user.LastName, user.Username, user.Password, user.AccountCreated, user.AccountUpdated)
	if err != nil {
		log.Printf(err.Error())
		return false
	}

	return true
}
