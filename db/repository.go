package db

import (
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib" // Standard library bindings for pgx (Postgres driver)
	"github.com/jmoiron/sqlx"
)

type User struct {
	UserId     int64  `gorm:"primaryKey autoIncrement"`
	ExternalId string `gorm:"uniqueIndex"`
}

type Repository struct {
	Db sqlx.DB
}

func NewRepository() (*Repository, error) {
	db, err := connect()
	if err != nil {
		return nil, err
	}

	return &Repository{Db: *db}, nil
}

func connect() (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", "postgres://invisibleprogrammer:invisiblepassword@localhost:5432/invisibleurl-db?sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("couldn't ping the database: %v", err)
	}

	log.Printf("Connected to the database")

	return db, nil
}
