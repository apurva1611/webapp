package main

import (
	"context"
	"database/sql"
"fmt"
	"log"
	"reflect"
	"time"

	_"github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error
const (
	username = "root"
	password = "pass1234"
	hostname = "127.0.0.1:3306"
	dbname   = "webappdb"
)

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func createDb() {
	db, err := sql.Open("mysql", dsn(""))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return
	}
	defer db.Close()

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

	db.Close()
	db, err = sql.Open("mysql", dsn(dbname))
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return
	}
	defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return
	}
	log.Printf("Connected to DB %s successfully\n", dbname)
}

func createTable(){
	_, err = db.Exec("DROP TABLE IF EXISTS mysql.users")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE mydb.users(
		id varchar(100) NOT NULL,
		firstname varchar(100) COLLATE utf8_unicode_ci NOT NULL,
		lastname varchar(100) COLLATE utf8_unicode_ci NOT NULL,
		email varchar(100) COLLATE utf8_unicode_ci NOT NULL,
		password varchar(255) COLLATE utf8_unicode_ci NOT NULL,
		created datetime NOT NULL,
		modified datetime NOT NULL,
		PRIMARY KEY (id))ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;`)
	if err != nil {
		panic(err)
	}

	// need to add argument
	_, err = db.Exec("INSERT INTO mysql.users (string_value, string_value, string_value, string_value, string_value, datetime_value, datetime_value ) VALUES (?, ?, ?, ?, ?, ?, ?)", "")
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT * FROM mysql.users")
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