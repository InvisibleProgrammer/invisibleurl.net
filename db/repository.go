package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	UserId     int64  `gorm:"primaryKey autoIncrement"`
	ExternalId string `gorm:"uniqueIndex"`
}

func Start() error {
	dsn := "host=localhost user=invisibleprogrammer password=invisiblepassword dbname=invisibleurl-db port=5432 sslmode=disable TimeZone=Europe/Budapest"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Couldn't connect to database: %v", err)
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("Db migration failed: %v", err)
	}

	return nil
}
