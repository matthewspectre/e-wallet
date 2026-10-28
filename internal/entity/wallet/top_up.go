package wallet

import "time"

const TopUpPaymentExpiry = 15 * time.Minute

type TopUp struct {
	ID          int64
	UserID      int64
	WalletID    int64
	OrderID     string
	Amount      int64
	Status      string
	SnapToken   string
	RedirectURL string
	PaymentType string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PaymentStatus struct {
	OrderID           string
	GrossAmount       string
	Currency          string
	TransactionStatus string
	FraudStatus       string
	PaymentType       string
}
