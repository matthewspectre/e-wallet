package user

import (
	entity "e-wallet/internal/entity/user"
	repo "e-wallet/internal/repository/user"

	"gorm.io/gorm"
)

type RepositoryMySQL struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repo.Repository {
	return &RepositoryMySQL{db: db}
}

func (r *RepositoryMySQL) Create(data *entity.User) error {
	return r.db.Table("user").Create(data).Error
}
