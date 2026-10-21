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

func (r *RepositoryMySQL) GetByID(id int64) (*entity.User, error) {
	var data entity.User
	if err := r.db.Table("user").First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *RepositoryMySQL) GetAll() ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64

	query := r.db.Table("user").Where("visible = ?", 1)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *RepositoryMySQL) HideByID(id int64) error {
	result := r.db.Table("user").Where("id = ?", id).Update("visible", 0)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		return nil
	}

	var count int64
	if err := r.db.Table("user").Where("id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
