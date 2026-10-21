package user

import "time"

type UserCreateRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	NoTelp   string `json:"no_telp"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	NoTelp    string    `json:"no_telp"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
