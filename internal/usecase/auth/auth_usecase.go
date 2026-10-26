package auth

import (
	"errors"
	"fmt"
	"strings"

	entity "e-wallet/internal/entity/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type UserRepository interface {
	GetByEmail(email string) (*entity.User, error)
}

type LoginResult struct {
	User *entity.User
}

type Usecase interface {
	Login(email, password string) (*LoginResult, error)
}

type usecase struct {
	users UserRepository
}

func NewUsecase(users UserRepository) Usecase {
	return &usecase{users: users}
}

func (u *usecase) Login(email, password string) (*LoginResult, error) {
	user, err := u.users.GetByEmail(strings.TrimSpace(email))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, fmt.Errorf("find login user: %w", err)
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, ErrInvalidCredentials
	}

	user.PasswordHash = ""
	return &LoginResult{
		User: user,
	}, nil
}
