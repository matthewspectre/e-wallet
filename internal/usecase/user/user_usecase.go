package user

import (
	entity "e-wallet/internal/entity/user"
	repo "e-wallet/internal/repository/user"
)

type Usecase interface {
	Create(data *entity.User) error
	GetByID(id int64) (*entity.User, error)
	GetAll() ([]*entity.User, int64, error)
	HideByID(id int64) error
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

func (u *usecase) GetByID(id int64) (*entity.User, error) {
	return u.repo.GetByID(id)
}

func (u *usecase) GetAll() ([]*entity.User, int64, error) {
	return u.repo.GetAll()
}

func (u *usecase) HideByID(id int64) error {
	return u.repo.HideByID(id)
}
