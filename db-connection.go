//https://www.bytesizego.com/blog/guide-to-logging-in-go

package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type attendeesdata struct {
	First_name  string `json:firstname`
	Second_name string `json:secondname`
	Age         uint   `json:age`
	Email       string `json:email`
	Occupation  string `json:occupation`
}

var jsonHandler = slog.NewJSONHandler(os.Stderr, nil)
var logger = slog.New(jsonHandler)

// logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

var db *sql.DB

func dbconnection() {

	cfg := mysql.Config{
		User:   "otieno",
		Passwd: "emmanuel.43",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "Booking",
	}

	var err error

	db, err = sql.Open("mysql", cfg.FormatDSN()) //Call sql.Open to initialize the db variable and connection to the db, passing the return value of FormatDSN.
	if err != nil {                              //Check for an error from sql.Open.
		log.Fatal(err) //log.Fatal will end execution and print the error to the console
	}

	pingErr := db.Ping() //Call db.Ping to confirm that connecting to the database works.
	if pingErr != nil {  //Check for an error from Ping, in case the connection failed
		log.Fatal(pingErr)
	}

	// err = db.ping(); if err != nil{ //checks whether theres an error with db.ping. if theres an error it exits
	// 	log.Fatal(err)  the output of db.ping will be passed assigned to err. if the output is not nil, the execution exits and the value assigned to err is printed
	// }

	log.Println("Connected!")
}

func Attendeescreation() {
	creation := `CREATE TABLE IF NOT EXISTS attendeesinfo (
	userid   INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	fitst_name    VARCHAR(255) NOT NULL,
	second_name   VARCHAR(255) NOT NULL,
	age      INT NOT NULL,
	email VARCHAR(255) NOT NULL UNIQUE,
	occupation VARCHAR(255) NOT NULL
	)`

	_, err := db.Exec(creation)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("created table")

}

func main() {
	dbconnection()
	Attendeescreation()
	router := gin.Default()
	router.POST("con/register", register)
	router.Run(":8088")
}

func register(c *gin.Context) {

	fmt.Println("kindly input your personal details")

	var userdata attendeesdata
	c.BindJSON(&userdata)

	query := `insert into attendeesinfo (fitst_name, second_name, age, email, occupation) values (?,?,?,?,?)`
	start := time.Now()
	_, err := db.Exec(query, userdata.First_name, userdata.Second_name, userdata.Age, userdata.Email, userdata.Occupation)
	elapsed := time.Since(start)
	if err != nil {
		logger.Error("ERROR Incurred", "error", err, "Execution time", elapsed)
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, userdata)
	logger.Info("user created", "email", userdata.Email, "Execution time", elapsed)

}
