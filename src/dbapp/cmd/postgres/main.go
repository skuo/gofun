package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const (
	DB_USER     = "devuser"
	DB_PASSWORD = "devpassword"
	DB_NAME     = "devdb"
)

func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	// ping
	err = db.Ping()
	checkErr(err)

	// insert
	fmt.Println("# Inserting values")
	var lastInsertID int
	err = db.QueryRow("INSERT INTO userinfo(username,departname,created) VALUES($1,$2,$3) returning uid;", "astaxie", "研发部门", "2012-12-09").
		Scan(&lastInsertID)
	checkErr(err)
	fmt.Println("last inserted id =", lastInsertID)

	// update
	fmt.Println("# Updating")
	stmt, err := db.Prepare("update userinfo set username=$1 where uid=$2")
	checkErr(err)

	res, err := stmt.Exec("astaxieupdate", lastInsertID)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect, "rows changed")

	// query
	fmt.Println("# Querying")
	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var uid int
		var username string
		var department string
		var created time.Time
		err = rows.Scan(&uid, &username, &department, &created)
		checkErr(err)
		fmt.Println("uid | username | department | created ")
		fmt.Printf("%3v | %8v | %6v | %6v\n", uid, username, department, created)
	}
	// check error at the end of for rows.next() to avoid calling rows.Close() inducing a runtime panic
	err = rows.Err()
	checkErr(err)

	// delete
	fmt.Println("# Deleting")
	stmt, err = db.Prepare("delete from userinfo where uid=$1")
	checkErr(err)

	res, err = stmt.Exec(lastInsertID)
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println(affect, "rows changed")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
