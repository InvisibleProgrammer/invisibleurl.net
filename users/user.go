package users

import "time"

type User struct {
	Id           int64
	PublicId     string
	EmailAddress string
	Activated    bool
	Status       int8
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
