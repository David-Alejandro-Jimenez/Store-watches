package database

import (
	"database/sql"
	"log"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() error {
	var err error

	DB, err = sql.Open("mysql", config.Config.Database.DSN)
	if err != nil {
		log.Printf("could not connect to database: %v", err)
		return err
	}

	err = DB.Ping()
	if err != nil {
		log.Printf("connection could not be verified: %v", err)
		return err
	}

	log.Printf("Database connection successful to %s@%s:%s/%s",
		config.Config.Database.User,
		config.Config.Database.Host,
		config.Config.Database.Port,
		config.Config.Database.Name)
	return nil
}
