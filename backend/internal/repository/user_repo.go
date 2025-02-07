package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/David-Alejandro-Jimenez/venta-relojes/internal/repository/database"
	"github.com/David-Alejandro-Jimenez/venta-relojes/pkg/security"
)

func GetUser(username string) (bool, error) {
	var existingUser bool
	var err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM User_Registration WHERE UserName = ?)", username).Scan(&existingUser)
	
		if err != nil {
			log.Println(err)
			return false, fmt.Errorf("error consultando la base de datos %w", err)
		}
	return existingUser, nil
}

func GetHashPassword(username string) (string, error) {
	var hashPassword string
	var query = "SELECT Password FROM User_Registration WHERE UserName= ?"
	var err = database.DB.QueryRow(query, username).Scan(&hashPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", err
		}
		return "", err
	}

	return hashPassword, nil
}	

func SaveUser(userName, password string) error {
	var salt, errSalt = security.GenerateSalt()
	if errSalt != nil {
		return errSalt
	}
	var hash, err = security.HashPassword(password, salt)
	if err != nil {
		return err
	}

	_, err = database.DB.Exec("INSERT INTO User_Registration (username, password, salt) VALUES (?, ?, ?)", userName, hash, salt)
	if err != nil {
		return err
	}
    return nil
}

func GetSalt(username string) (string, error) {
	var salt string
	var query = "SELECT Salt FROM User_Registration WHERE UserName= ?"
	var err  = database.DB.QueryRow(query, username).Scan(&salt)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", err
		}
		return "", err
	}

	return salt, nil
}
