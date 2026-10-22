package wallet

import "time"

type WalletModel struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;column:id;type:bigint"`
	UserID    int64     `gorm:"column:user_id;uniqueIndex;not null"`
	Balance   int64     `gorm:"column:balance;not null;default:0"`
	Currency  string    `gorm:"column:currency;size:3;not null;default:IDR"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (WalletModel) TableName() string {
	return "wallets"
}
