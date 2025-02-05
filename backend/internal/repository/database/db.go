package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/spf13/viper"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() error {
	var err error
	var user = viper.GetString("DB_USER")
	var password = viper.GetString("DB_PASSWORD")
	var host = viper.GetString("DB_HOST")
	var port = viper.GetString("DB_PORT")
	var database = viper.GetString("DB_NAME")

	var dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4", user, password, host, port, database)
	
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("could not connect to database: %v", err)
		return err
	}

	err = DB.Ping() 
	if err != nil {
		log.Printf("connection could not be verified: %v", err)
        return err
    }

	log.Println("Database connection successful")
	return nil
}