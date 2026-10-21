package user

import "time"

type User struct {
	ID           int64
	Name         string
	Email        string
	NoTelp       string
	PasswordHash string
	Visible      int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
