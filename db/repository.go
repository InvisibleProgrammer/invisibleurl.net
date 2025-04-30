package db

import (
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib" // Standard library bindings for pgx (Postgres driver)
	"github.com/jmoiron/sqlx"
	"invisibleprogrammer.com/invisibleurl/environment"
)

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

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		environment.DB_USER,
		environment.DB_PASSWORD,
		environment.DB_HOST,
		environment.DB_PORT,
		environment.DB_NAME,
	)

	log.Printf("DB host: %s\n", environment.DB_HOST)

	db, err := sqlx.Open("pgx", connString)
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
