package wallet

import (
	entity "e-wallet/internal/entity/wallet"
	model "e-wallet/internal/model/wallet"
	repo "e-wallet/internal/repository/wallet"
	"math"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *RepositoryMySQL) GetTopUpsByUserID(userID int64) ([]*entity.TopUp, error) {
	type topUpWithUserID struct {
		ID          int64     `gorm:"column:id"`
		UserID      int64     `gorm:"column:user_id"`
		WalletID    int64     `gorm:"column:wallet_id"`
		OrderID     string    `gorm:"column:order_id"`
		Amount      int64     `gorm:"column:amount"`
		Status      string    `gorm:"column:status"`
		SnapToken   string    `gorm:"column:snap_token"`
		RedirectURL string    `gorm:"column:redirect_url"`
		PaymentType string    `gorm:"column:payment_type"`
		CreatedAt   time.Time `gorm:"column:created_at"`
		UpdatedAt   time.Time `gorm:"column:updated_at"`
	}

	expiredBefore := time.Now().Add(-entity.TopUpPaymentExpiry)
	if err := r.db.Model(&model.TopUpModel{}).
		Where("status = ? AND created_at <= ? AND wallet_id IN (?)", "pending", expiredBefore,
			r.db.Table("wallets").Select("id").Where("user_id = ?", userID)).
		Updates(map[string]interface{}{
			"status":     "expired",
			"updated_at": gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error; err != nil {
		return nil, err
	}

	var rows []topUpWithUserID
	err := r.db.Table("wallet_topups AS wt").
		Select("wt.id, w.user_id, wt.wallet_id, wt.order_id, wt.amount, wt.status, wt.snap_token, wt.redirect_url, wt.payment_type, wt.created_at, wt.updated_at").
		Joins("JOIN wallets AS w ON w.id = wt.wallet_id").
		Where("w.user_id = ?", userID).
		Order("wt.created_at DESC, wt.id DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	topUps := make([]*entity.TopUp, 0, len(rows))
	for _, row := range rows {
		topUps = append(topUps, &entity.TopUp{
			ID:          row.ID,
			UserID:      row.UserID,
			WalletID:    row.WalletID,
			OrderID:     row.OrderID,
			Amount:      row.Amount,
			Status:      row.Status,
			SnapToken:   row.SnapToken,
			RedirectURL: row.RedirectURL,
			PaymentType: row.PaymentType,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
		})
	}
	return topUps, nil
}

func (r *RepositoryMySQL) CreateTopUp(userID int64, data *entity.TopUp) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var wallet model.WalletModel
		if err := tx.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			return err
		}

		data.WalletID = wallet.ID
		topUp := model.TopUpModel{
			WalletID:  data.WalletID,
			OrderID:   data.OrderID,
			Amount:    data.Amount,
			Status:    "pending",
			CreatedAt: data.CreatedAt,
			UpdatedAt: data.UpdatedAt,
		}
		if err := tx.Create(&topUp).Error; err != nil {
			return err
		}

		data.ID = topUp.ID
		data.Status = topUp.Status
		return nil
	})
}

func (r *RepositoryMySQL) UpdateTopUpPayment(orderID, token, redirectURL string) error {
	return r.db.Model(&model.TopUpModel{}).
		Where("order_id = ?", orderID).
		Updates(map[string]interface{}{
			"snap_token":   token,
			"redirect_url": redirectURL,
			"updated_at":   gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error
}

func (r *RepositoryMySQL) GetTopUpByOrderID(orderID string) (*entity.TopUp, error) {
	var data model.TopUpModel
	if err := r.db.Where("order_id = ?", orderID).First(&data).Error; err != nil {
		return nil, err
	}
	return toTopUpEntity(&data), nil
}

func (r *RepositoryMySQL) ApplyTopUpStatus(orderID, status, paymentType string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var topUp model.TopUpModel
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("order_id = ?", orderID).
			First(&topUp).Error; err != nil {
			return err
		}

		if topUp.Status == "paid" || (topUp.Status != "pending" && status == "pending") {
			return nil
		}

		if status == "paid" {
			var wallet model.WalletModel
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				First(&wallet, topUp.WalletID).Error; err != nil {
				return err
			}
			if topUp.Amount > 0 && wallet.Balance > math.MaxInt64-topUp.Amount {
				return gorm.ErrInvalidData
			}
			wallet.Balance += topUp.Amount
			if err := tx.Save(&wallet).Error; err != nil {
				return err
			}
		}

		return tx.Model(&topUp).Updates(map[string]interface{}{
			"status":       status,
			"payment_type": paymentType,
			"updated_at":   gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error
	})
}

func toTopUpEntity(data *model.TopUpModel) *entity.TopUp {
	return &entity.TopUp{
		ID:          data.ID,
		WalletID:    data.WalletID,
		OrderID:     data.OrderID,
		Amount:      data.Amount,
		Status:      data.Status,
		SnapToken:   data.SnapToken,
		RedirectURL: data.RedirectURL,
		PaymentType: data.PaymentType,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
}
