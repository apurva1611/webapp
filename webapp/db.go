package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

var db *sql.DB

// const (
// 	username = "root"
// 	password = "Rajuabha25!"
// 	hostname = "localhost:3306"
// 	dbname   = "webappdb"
// )

// func dsn() string {
// 	//dsurl := localhost
// 	//hostname := rdsurl + port
// 	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)
// }
const (
	username = "adminuser@poller-nstance"
	password = "Pass1234"
	// port     = ":3306"
	dbname = "webappdb"
)

func dsn() string {
	rdsurl := os.Getenv("rdsurl")
	hostname := rdsurl
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)
}

func openDB() {
	var err error
	db, err = sql.Open("mysql", dsn())
	if err != nil {
		log.Error("Error %s when opening DB\n", err)
		panic(err)
	}
}

func closeDB() {
	db.Close()
}

func dbHealthCheck() error {
	err := db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func createDb() {
	openDB()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)

	if err != nil {
		log.Error("Error %s when creating DB\n", err)
		return
	}

	no, err := res.RowsAffected()
	if err != nil {
		log.Error("Error %s when fetching rows", err)
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
		log.Error("Errors %s pinging DB", err)
		return
	}
	log.Info("Connected to DB %s successfully\n", dbname)
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
		value float NOT NULL,
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
	//timer := prometheus.NewTimer(requestDuration)
	//defer timer.ObserveDuration()
	user := User{}
	err := db.QueryRow(`SELECT id, firstname, lastname, username, created, modified 
							FROM webappdb.users WHERE id = ?`, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username,
		&user.AccountCreated, &user.AccountUpdated)
	//time.Sleep(time.Duration(rand.NormFloat64()*10000+50000) * time.Microsecond)
	if err != nil {
		log.Error("User query by id failed")
		log.Error(err.Error())
		return nil
	}
	log.Info("User query by id succeeded")
	return &user
}
func queryById(id string) *User {
	user := User{}
	err := db.QueryRow(`SELECT id, firstname, lastname, username, created, modified 
							FROM webappdb.users WHERE id = ?`, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username,
		&user.AccountCreated, &user.AccountUpdated)
	if err != nil {
		log.Error("User query by id failed")
		log.Error(err.Error())
		return nil
	}
	log.Info("User query by id succeeded")
	return &user
}

func queryByUsername(username string) *User {
	user := User{}
	err := db.QueryRow(`SELECT id, firstname, lastname, username, created, modified 
							FROM webappdb.users WHERE username = ?`, username).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username,
		&user.AccountCreated, &user.AccountUpdated)
	if err != nil {
		log.Error("User query by username failed")
		log.Error(err.Error())
		return nil
	}
	log.Info("User query by username succeeded")
	return &user
}

func insertUser(user User) bool {
	insert, err := db.Prepare(`INSERT INTO webappdb.users(id, firstname, lastname, username, password, created, modified) 
						VALUES (?, ?, ?, ?, ?, ?, ?)`)

	if err != nil {
		log.Error("Insert user query error")
		log.Error(err.Error())
		return false
	}

	_, err = insert.Exec(user.ID, user.FirstName, user.LastName, user.Username, user.Password, user.AccountCreated, user.AccountUpdated)
	if err != nil {
		log.Error("Insert user query error")
		log.Error(err.Error())
		return false
	}
	log.Info("Insert user query succeeded")
	return true
}

func updateUser(user User) bool {
	update, err := db.Prepare(`UPDATE webappdb.users SET firstname=?, lastname=?, password=?, modified=? 
										WHERE id=?`)

	if err != nil {
		log.Error("Update user query error")
		log.Error(err.Error())
		return false
	}

	_, err = update.Exec(user.FirstName, user.LastName, user.Password, user.AccountUpdated, user.ID)
	if err != nil {
		log.Error("Update user query error")
		log.Error(err.Error())
		return false
	}
	log.Info("Update user query succeeded")
	return true
}

func insertWatch(watch WATCH) bool {
	insert, err := db.Prepare(`INSERT INTO webappdb.watch(watch_id, user_id,zipcode, alerts, watch_created, watch_updated) 
						VALUES (?, ?, ?, ?, ?, ?)`)

	if err != nil {
		log.Error("Insert watch query error")
		log.Error(err.Error())
		return false
	}
	alerts_json, err := json.Marshal(&watch.Alerts)
	_, err = insert.Exec(watch.ID, watch.UserId, watch.Zipcode, alerts_json, watch.WatchCreated, watch.WatchUpdated)
	if err != nil {
		log.Error("Insert watch query error")
		log.Error(err.Error())
		return false
	}
	log.Info("Insert watch query succeeded")
	return true
}
func queryWatchByUserId(id string) *[]WATCH {
	var watches []WATCH
	rows, err := db.Query(`SELECT watch_id, user_id, zipcode,watch_created, watch_updated 
							FROM webappdb.watch WHERE user_id = ?`, id)

	defer rows.Close()
	for rows.Next() {
		watch := WATCH{}
		var alerts []ALERT
		err = rows.Scan(&watch.ID, &watch.UserId, &watch.Zipcode, &watch.WatchCreated, &watch.WatchUpdated)
		alerts_received := queryAlertsByWatchId(watch.ID)
		//fmt.println(string(*alerts_received))
		watch.Alerts = alerts
		for _, element := range *alerts_received {
			watch.Alerts = append(watch.Alerts, element)
		}
		//fmt.println(watch.Alerts[0].ID)
		if err != nil {
			// handle this error
			log.Error("Get watch query by user id failed")
			panic(err)
		}
		watches = append(watches, watch)

	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		log.Error("Get watch query by user id failed")
		log.Error(err.Error())
		return nil
	}
	log.Info("Get watch by user ID query succeeded")
	return &watches
}

func queryByWatchID(id string) *WATCH {
	watch := WATCH{}
	err := db.QueryRow(`SELECT watch_id, user_id, zipcode, watch_created,watch_updated
							FROM webappdb.watch WHERE watch_id = ?`, id).Scan(&watch.ID, &watch.UserId, &watch.Zipcode, &watch.WatchCreated, &watch.WatchUpdated)
	if err != nil {
		log.Error("Get watch query by watch id failed")
		log.Error(err.Error())
		return nil
	}
	var alerts []ALERT
	alerts_received := queryAlertsByWatchId(id)
	watch.Alerts = alerts
	for _, element := range *alerts_received {
		watch.Alerts = append(watch.Alerts, element)
	}
	log.Info("Get watch query by watch id succeeded")
	return &watch

}

func insertAlert(alert ALERT) bool {
	insert, err := db.Prepare(`INSERT INTO webappdb.alert(alert_id, watch_id, field_type, operator, value, alert_created, alert_updated) 
						VALUES (?, ?, ?, ?, ?, ?, ?)`)

	if err != nil {
		log.Error("insert alert query failed")
		log.Error(err.Error())
		return false
	}
	_, err = insert.Exec(alert.ID, alert.WatchId, alert.FieldType, alert.Operator, alert.Value, alert.AlertCreated, alert.AlertUpdated)
	if err != nil {
		log.Error("insert alert query failed")
		log.Error(err.Error())
		return false
	}
	log.Info("insert alert query succeeded")
	return true
}

func queryAlertsByWatchId(id string) *[]ALERT {
	var alerts []ALERT
	rows, err := db.Query(`SELECT alert_id,field_type, operator, value,alert_created,alert_updated 
							FROM webappdb.alert WHERE watch_id = ?`, id)
	defer rows.Close()
	for rows.Next() {
		alert := ALERT{}
		err = rows.Scan(&alert.ID, &alert.FieldType, &alert.Operator, &alert.Value, &alert.AlertCreated, &alert.AlertUpdated)
		if err != nil {
			// handle this error
			log.Error("GETB query alerts by watch id error")
			panic(err)
		}
		alerts = append(alerts, alert)

	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		log.Error("GET query alerts by watch id error")
		log.Error(err.Error())
		return nil
	}
	log.Info("GET query alerts by watch id succeeded")
	return &alerts

}

func deleteAlert(id string) bool {
	delete, err := db.Prepare("DELETE FROM webappdb.alert WHERE alert_id=?")

	delete.Exec(id)
	if err != nil {
		log.Error("Delete alert by alert id error")
		log.Error(err.Error())
		return false
	}
	log.Info("Delete alert by alert id")
	return true
}
func updateWatch(watch WATCH) bool {
	update, err := db.Prepare(`UPDATE webappdb.watch SET watch_id=?, user_id=?, zipcode=?, alerts=? , watch_created=?, watch_updated=?
										WHERE watch_id=?`)

	if err != nil {
		log.Error("Update watch by watch id error")
		log.Error(err.Error())
		return false
	}
	alerts_json, err := json.Marshal(&watch.Alerts)
	_, err = update.Exec(watch.ID, watch.UserId, watch.Zipcode, alerts_json, watch.WatchCreated, watch.WatchUpdated, watch.ID)
	if err != nil {
		log.Error("Update watch by watch id error")
		log.Error(err.Error())
		return false
	}
	log.Info("Update watch by watch id query succeeded")
	return true
}

func deleteWatch(id string) bool {
	delete, err := db.Prepare("DELETE FROM webappdb.watch WHERE watch_id=?")

	delete.Exec(id)
	if err != nil {
		log.Error("Delete watch by watch id error")
		log.Error(err.Error())
		return false
	}
	log.Info("Delete watch by watch id succeeded")
	return true
}
