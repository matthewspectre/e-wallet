package wallet

import "e-wallet/internal/entity/wallet"

type Repository interface {
	GetByUserID(userID int64) (*wallet.Wallet, error)
}
