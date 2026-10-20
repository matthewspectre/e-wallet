package user

import "time"

type UserModel struct {
	ID           int64     `gorm:"primaryKey;autoIncrement;column:id;type:bigint"`
	Name         string    `gorm:"column:name"`
	Email        string    `gorm:"column:email;unique"`
	PasswordHash string    `gorm:"column:password_hash"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (UserModel) TableName() string {
	return "user"
}
