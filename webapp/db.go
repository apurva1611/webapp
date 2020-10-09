package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

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

	ctx, cancelFunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return
	}
	log.Printf("Connected to DB %s successfully\n", dbname)
}

func createTable() {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS webappdb.users(
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

func queryByUsername(username string) *User {
	user := User{}
	err := db.QueryRow(`SELECT id, firstname, lastname, username, created, modified 
							FROM webappdb.users WHERE username = ?`, username).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username,
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

func updateUser(user User) bool {
	update, err := db.Prepare(`UPDATE webappdb.users SET firstname=?, lastname=?, username=?, password=?, modified=? 
										WHERE id=?`)

	if err != nil {
		log.Printf(err.Error())
		return false
	}

	_, err = update.Exec(user.FirstName, user.LastName, user.Username, user.Password, user.AccountUpdated, user.ID)
	if err != nil {
		log.Printf(err.Error())
		return false
	}

	return true
}
