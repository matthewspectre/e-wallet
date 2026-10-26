package user

import (
	"e-wallet/internal/entity/user"
)

type Repository interface {
	Create(data *user.User) error
	GetByID(id int64) (*user.User, error)
	GetByEmail(email string) (*user.User, error)
	GetAll() ([]*user.User, int64, error)
	HideByID(id int64) error
}
