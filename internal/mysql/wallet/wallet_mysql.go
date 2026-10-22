package wallet

import (
	entity "e-wallet/internal/entity/wallet"
	model "e-wallet/internal/model/wallet"
	repo "e-wallet/internal/repository/wallet"

	"gorm.io/gorm"
)

type RepositoryMySQL struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repo.Repository {
	return &RepositoryMySQL{db: db}
}

func (r *RepositoryMySQL) GetByUserID(userID int64) (*entity.Wallet, error) {
	var data model.WalletModel
	if err := r.db.Where("user_id = ?", userID).First(&data).Error; err != nil {
		return nil, err
	}

	return &entity.Wallet{
		ID:        data.ID,
		UserID:    data.UserID,
		Balance:   data.Balance,
		Currency:  data.Currency,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}, nil
}
