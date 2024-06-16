package users

import "time"

type UserPassword struct {
	UserId       int64
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
