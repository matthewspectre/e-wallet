package wallet

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	entity "e-wallet/internal/entity/wallet"
	repo "e-wallet/internal/repository/wallet"
)

var (
	ErrInvalidNotification = errors.New("invalid Midtrans notification")
	ErrInvalidTopUpAmount  = errors.New("top-up amount must be greater than zero")
	ErrUnexpectedPayment   = errors.New("Midtrans transaction does not match top-up")
)

type PaymentGateway interface {
	CreateSnapTransaction(orderID string, amount int64) (token, redirectURL string, err error)
	VerifyNotification(orderID, statusCode, grossAmount, signature string) bool
	GetTransactionStatus(orderID string) (*entity.PaymentStatus, error)
}

type Usecase interface {
	GetByUserID(userID int64) (*entity.Wallet, error)
	CreateTopUp(userID, amount int64) (*entity.TopUp, error)
	HandleNotification(orderID, statusCode, grossAmount, signature string) error
}

type usecase struct {
	repo    repo.Repository
	payment PaymentGateway
}

func NewUsecase(r repo.Repository, payment PaymentGateway) Usecase {
	return &usecase{repo: r, payment: payment}
}

func (u *usecase) GetByUserID(userID int64) (*entity.Wallet, error) {
	return u.repo.GetByUserID(userID)
}

func (u *usecase) CreateTopUp(userID, amount int64) (*entity.TopUp, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("user ID must be positive")
	}
	if amount <= 0 {
		return nil, ErrInvalidTopUpAmount
	}
	if u.payment == nil {
		return nil, fmt.Errorf("Midtrans payment gateway is not configured")
	}

	orderID, err := newOrderID(userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	topUp := &entity.TopUp{
		OrderID:   orderID,
		Amount:    amount,
		Status:    "pending",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := u.repo.CreateTopUp(userID, topUp); err != nil {
		return nil, err
	}

	token, redirectURL, err := u.payment.CreateSnapTransaction(orderID, amount)
	if err != nil {
		if updateErr := u.repo.ApplyTopUpStatus(orderID, "failed", ""); updateErr != nil {
			return nil, fmt.Errorf("create Midtrans transaction: %w (also failed to mark top-up failed: %v)", err, updateErr)
		}
		return nil, err
	}
	if err := u.repo.UpdateTopUpPayment(orderID, token, redirectURL); err != nil {
		return nil, fmt.Errorf("save Midtrans payment details: %w", err)
	}

	topUp.SnapToken = token
	topUp.RedirectURL = redirectURL
	return topUp, nil
}

func (u *usecase) HandleNotification(orderID, statusCode, grossAmount, signature string) error {
	if u.payment == nil || !u.payment.VerifyNotification(orderID, statusCode, grossAmount, signature) {
		return ErrInvalidNotification
	}

	topUp, err := u.repo.GetTopUpByOrderID(orderID)
	if err != nil {
		return err
	}
	if !amountMatches(grossAmount, topUp.Amount) {
		return ErrUnexpectedPayment
	}

	status, err := u.payment.GetTransactionStatus(orderID)
	if err != nil {
		return err
	}
	if status.OrderID != orderID ||
		status.Currency != "IDR" ||
		!amountMatches(status.GrossAmount, topUp.Amount) {
		return ErrUnexpectedPayment
	}

	localStatus, err := mapPaymentStatus(status.TransactionStatus, status.FraudStatus)
	if err != nil {
		return err
	}
	return u.repo.ApplyTopUpStatus(orderID, localStatus, status.PaymentType)
}

func newOrderID(userID int64) (string, error) {
	var randomID [8]byte
	if _, err := rand.Read(randomID[:]); err != nil {
		return "", fmt.Errorf("generate top-up order ID: %w", err)
	}
	return fmt.Sprintf("EW-%d-%d-%s", userID, time.Now().UnixNano(), hex.EncodeToString(randomID[:])), nil
}

func amountMatches(amount string, expected int64) bool {
	actual, ok := new(big.Rat).SetString(strings.TrimSpace(amount))
	if !ok {
		return false
	}
	return actual.Cmp(new(big.Rat).SetInt64(expected)) == 0
}

func mapPaymentStatus(status, fraudStatus string) (string, error) {
	switch strings.ToLower(status) {
	case "settlement":
		return "paid", nil
	case "capture":
		switch strings.ToLower(fraudStatus) {
		case "accept":
			return "paid", nil
		case "challenge":
			return "pending", nil
		default:
			return "failed", nil
		}
	case "pending":
		return "pending", nil
	case "expire":
		return "expired", nil
	case "cancel":
		return "cancelled", nil
	case "deny":
		return "failed", nil
	default:
		return "", fmt.Errorf("%w: unsupported transaction status %q", ErrUnexpectedPayment, status)
	}
}
