package wallet

import (
	"time"

	"e-wallet/internal/entity/wallet"
)

type Repository interface {
	GetByUserID(userID int64) (*wallet.Wallet, error)
	GetTopUpsByUserID(userID int64, startDate, endDate *time.Time) ([]*wallet.TopUp, error)
	CreateTopUp(userID int64, data *wallet.TopUp) error
	UpdateTopUpPayment(orderID, token, redirectURL string) error
	GetTopUpByOrderID(orderID string) (*wallet.TopUp, error)
	ApplyTopUpStatus(orderID, status, paymentType string) error
}
