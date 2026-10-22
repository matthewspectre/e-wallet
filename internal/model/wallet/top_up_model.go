package wallet

import "time"

type TopUpModel struct {
	ID          int64     `gorm:"primaryKey;autoIncrement;column:id;type:bigint"`
	WalletID    int64     `gorm:"column:wallet_id;not null;index"`
	OrderID     string    `gorm:"column:order_id;size:64;not null;uniqueIndex"`
	Amount      int64     `gorm:"column:amount;not null"`
	Status      string    `gorm:"column:status;size:20;not null;default:pending"`
	SnapToken   string    `gorm:"column:snap_token;size:255"`
	RedirectURL string    `gorm:"column:redirect_url;size:1024"`
	PaymentType string    `gorm:"column:payment_type;size:50"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (TopUpModel) TableName() string {
	return "wallet_topups"
}
