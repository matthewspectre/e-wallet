package auth

import (
	"errors"
	"testing"

	entity "e-wallet/internal/entity/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userRepositoryStub struct {
	user *entity.User
	err  error
}

func (r userRepositoryStub) GetByEmail(string) (*entity.User, error) {
	return r.user, r.err
}

func TestLoginReturnsPublicUser(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("correct password"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	user := &entity.User{
		ID:           42,
		Name:         "Wallet User",
		Email:        "wallet@example.com",
		PasswordHash: string(passwordHash),
		Visible:      1,
	}

	result, err := NewUsecase(userRepositoryStub{user: user}).
		Login(" wallet@example.com ", "correct password")
	if err != nil {
		t.Fatalf("Login returned an error: %v", err)
	}
	if result.User.ID != 42 || result.User.Email != user.Email {
		t.Fatalf("unexpected login result: %#v", result)
	}
	if result.User.PasswordHash != "" {
		t.Fatal("password hash must not be included in the login result")
	}
}

func TestLoginRejectsInvalidCredentials(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("correct password"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		repo userRepositoryStub
		pass string
	}{
		{
			name: "wrong password",
			repo: userRepositoryStub{user: &entity.User{PasswordHash: string(passwordHash)}},
			pass: "wrong password",
		},
		{
			name: "unknown account",
			repo: userRepositoryStub{err: gorm.ErrRecordNotFound},
			pass: "any password",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewUsecase(tt.repo).Login("wallet@example.com", tt.pass)
			if !errors.Is(err, ErrInvalidCredentials) {
				t.Fatalf("expected ErrInvalidCredentials, got %v", err)
			}
		})
	}
}
