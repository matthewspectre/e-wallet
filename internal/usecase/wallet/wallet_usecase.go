package wallet

import (
	entity "e-wallet/internal/entity/wallet"
	repo "e-wallet/internal/repository/wallet"
)

type Usecase interface {
	GetByUserID(userID int64) (*entity.Wallet, error)
}

type usecase struct {
	repo repo.Repository
}

func NewUsecase(r repo.Repository) Usecase {
	return &usecase{repo: r}
}

func (u *usecase) GetByUserID(userID int64) (*entity.Wallet, error) {
	return u.repo.GetByUserID(userID)
}
