package authRepository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/security_auth"
)

type UserRepository interface {
    UserExists(username string) (bool, error)
    GetHashPassword(username string) (string, error)
    GetSalt(username string) (string, error)
    SaveUser(username, password string) error
}

type userRepository struct {
	db *sql.DB
    saltGenerator securityAuth.Generator
    hasher        securityAuth.Hasher
}

func NewUserRepository(db *sql.DB, saltGenerator securityAuth.Generator, hasher securityAuth.Hasher) UserRepository {
    if db == nil {
        log.Fatal("ERROR: database.DB is nil in NewUserRepository")
    }

    if saltGenerator == nil {
        log.Fatal("ERROR: saltGenerator is nil in NewUserRepository")
    }
    if hasher == nil {
        log.Fatal("ERROR: hasher is nil in NewUserRepository")
    }

    log.Println("NewUserRepository() is running successfully")

    return &userRepository{
        db: db,
        saltGenerator: saltGenerator,
        hasher: hasher,
    }
}

func (a *userRepository) UserExists(username string) (bool, error) {
    var exists bool
    query := "SELECT EXISTS(SELECT 1 FROM User_Registration WHERE UserName = ?)"
    err := a.db.QueryRow(query, username).Scan(	&exists)
    if err != nil {
        return false, fmt.Errorf("error querying the database: %w", err)
    }
    return exists, nil
}

func (a *userRepository) GetHashPassword(username string) (string, error) {
    var hashPassword string
    query := "SELECT Password FROM User_Registration WHERE UserName = ?"
    err := a.db.QueryRow(query, username).Scan(&hashPassword)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", err
        }
        return "", err
    }
    return hashPassword, nil
}

func (a *userRepository) GetSalt(username string) (string, error) {
    var salt string
    query := "SELECT Salt FROM User_Registration WHERE UserName = ?"
    err := a.db.QueryRow(query, username).Scan(&salt)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", err
        }
        return "", err
    }
    return salt, nil
}


func (a *userRepository) SaveUser(username, password string) error {
    salt, err := a.saltGenerator.Generate()
	if err != nil {
        log.Println("error save user 1", err)
		return err
	}

    combined := securityAuth.Combined(password, salt)
    hash, err := a.hasher.Hash(combined)
	if err != nil {
        log.Println("error save user 2", err)
		return err
	}

	_, err = a.db.Exec("INSERT INTO User_Registration (UserName, Password, Salt) VALUES (?, ?, ?)", username, hash, salt)
	if err != nil {
		return err
	}
    return nil
}
