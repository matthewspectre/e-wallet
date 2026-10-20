package user

import (
	"e-wallet/internal/entity/user"
)

type Repository interface {
	Create(data *user.User) error
}
