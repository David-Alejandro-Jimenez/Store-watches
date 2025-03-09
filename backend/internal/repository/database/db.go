package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var DB *sql.DB

// The InitDB function is responsible for initializing the connection to the database using configuration parameters obtained with Viper.
// 1. Configuration: Get connection data using Viper.
// 2. DSN: Constructs a connection string for MySQL.
// 3. Connection: Open and verify the connection to the database.
// 4. Errors: These are handled and returned in case of problems, ensuring that the application knows if the connection was successful or not.
// This function is essential to establish and verify the connection to the database before executing operations that depend on it.
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