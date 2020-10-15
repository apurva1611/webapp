package main

import (
	"encoding/json"
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
	hostname = "localhost:3306"
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
		username varchar(100) COLLATE utf8_unicode_ci NOT NULL UNIQUE,
		password varchar(255) COLLATE utf8_unicode_ci NOT NULL,
		created datetime NOT NULL,
		modified datetime NOT NULL,
		PRIMARY KEY (id, username))ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;`)
    // create table watch
	_, err1 := db.Exec(`CREATE TABLE IF NOT EXISTS webappdb.watch(
		watch_id varchar(100) NOT NULL,
		user_id varchar(100) COLLATE utf8_unicode_ci NOT NULL,
		zipcode varchar(100) COLLATE utf8_unicode_ci NOT NULL,
		alerts json NOT NULL,
		watch_created datetime NOT NULL,
		watch_updated datetime NOT NULL,
		PRIMARY KEY (watch_id),
		FOREIGN KEY (user_id) REFERENCES users(id) 
		)ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;`)
    // create table alert
	_, err2 := db.Exec(`CREATE TABLE IF NOT EXISTS webappdb.alert(
		alert_id varchar(100) NOT NULL,
		watch_id varchar(100) COLLATE utf8_unicode_ci NOT NULL,
		field_type ENUM('temp', 'feels_like', 'temp_min', 'temp_max', 'pressure','humidity') COLLATE utf8_unicode_ci NOT NULL,
		operator ENUM('gt', 'gte', 'eq', 'lt', 'lte') COLLATE utf8_unicode_ci NOT NULL,
		value int NOT NULL,
		alert_created datetime NOT NULL,
		alert_updated datetime NOT NULL,
		PRIMARY KEY (alert_id),
		FOREIGN KEY (watch_id) REFERENCES watch(watch_id) 
		)ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;`)
	
	if err != nil {
		panic(err)
	}
	if err1 != nil {
		panic(err1)
	}
	if err2 != nil {
		panic(err2)
	}
}

func queryByID(id string) *User {
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
	update, err := db.Prepare(`UPDATE webappdb.users SET firstname=?, lastname=?, password=?, modified=? 
										WHERE id=?`)

	if err != nil {
		log.Printf(err.Error())
		return false
	}

	_, err = update.Exec(user.FirstName, user.LastName, user.Password, user.AccountUpdated, user.ID)
	if err != nil {
		log.Printf(err.Error())
		return false
	}

	return true
}

func insertWatch(watch WATCH)bool {
	insert, err := db.Prepare(`INSERT INTO webappdb.watch(watch_id, user_id,zipcode, alerts, watch_created, watch_updated) 
						VALUES (?, ?, ?, ?, ?, ?)`)

	if err != nil {
		log.Printf(err.Error())
		return false
	}
	alerts_json,err := json.Marshal(&watch.Alerts)
	fmt.Println(string(alerts_json))
	_, err = insert.Exec(watch.ID, watch.UserId, watch.Zipcode,alerts_json, watch.WatchCreated,watch.WatchUpdated)
	if err != nil {
		log.Printf(err.Error())
		return false
	}
	fmt.Println("Watch inserted")
	return true
}
// func queryWatchByUserId(id string) *WATCHES {
// 	watches = WATCHES
// 	rows,err := db.Query(`SELECT watch_id, user_id, zipcode, alerts, watch_created, watch_updated 
// 							FROM webappdb.watch WHERE id = ?`, id)
// 	defer rows.Close()
//     var int i
// 	for rows.Next() {
// 		watch := WATCH{}
// 		err = rows.Scan(&watch.ID, &watch.UserId,&watch.Zipcode, &watch.Alerts,&watch.WatchCreated, &watch.WatchUpdated)
// 		if err != nil {
// 		// handle this error
// 		panic(err)
// 		}
// 		//fmt.Println(id, firstName)
// 		watches[i]=watch
// 		i++
		
// 	}
// 	// get any error encountered during iteration
// 	err = rows.Err()
// 	if err != nil {
// 		panic(err)
// 		return nil
// 	}

// 	return &watches
// }

func queryByWatchID(id string) *WATCH {
	watch := WATCH{}
	err := db.QueryRow(`SELECT watch_id, user_id, zipcode, alerts, watch_created, watch_updated 
							FROM webappdb.watch WHERE watch_id = ?`, id).Scan(&watch.ID, &watch.UserId, &watch.Zipcode, &watch.Alerts,&watch.WatchCreated, &watch.WatchUpdated)
	if err != nil {
		log.Printf(err.Error())
		return nil
	}
	//fmt.printf(watch)
	return &watch
}

func insertAlert(alert ALERT) bool {
	insert, err := db.Prepare(`INSERT INTO webappdb.alert(alert_id, watch_id, field_type, operator, value, alert_created, alert_updated) 
						VALUES (?, ?, ?, ?, ?, ?, ?)`)

	if err != nil {
		log.Printf(err.Error())
		return false
	}
	_, err = insert.Exec(alert.ID, alert.WatchId, alert.FieldType, alert.Operator, alert.Value, alert.AlertCreated, alert.AlertUpdated)
	if err != nil {
		log.Printf(err.Error())
		return false
	}

	return true
}