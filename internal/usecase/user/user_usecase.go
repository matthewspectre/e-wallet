package user

import (
	entity "e-wallet/internal/entity/user"
	repo "e-wallet/internal/repository/user"
)

type Usecase interface {
	Create(data *entity.User) error
}

type usecase struct {
	repo repo.Repository
}

func NewUsecase(r repo.Repository) Usecase {
	return &usecase{repo: r}
}

func (u *usecase) Create(data *entity.User) error {
	return u.repo.Create(data)
}
