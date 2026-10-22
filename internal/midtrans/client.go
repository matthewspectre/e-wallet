package midtrans

import (
	"crypto/sha512"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strings"

	entity "e-wallet/internal/entity/wallet"
	midtransSDK "github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type Client struct {
	serverKey string
	snap      snap.Client
	core      coreapi.Client
}

func NewClient(serverKey string) *Client {
	serverKey = strings.TrimSpace(serverKey)
	snapClient := snap.Client{}
	snapClient.New(serverKey, midtransSDK.Sandbox)
	coreClient := coreapi.Client{}
	coreClient.New(serverKey, midtransSDK.Sandbox)

	return &Client{
		serverKey: serverKey,
		snap:      snapClient,
		core:      coreClient,
	}
}

func (c *Client) CreateSnapTransaction(orderID string, amount int64) (string, string, error) {
	if c.serverKey == "" {
		return "", "", fmt.Errorf("MIDTRANS_SERVER_KEY is not configured")
	}

	items := []midtransSDK.ItemDetails{{
		ID:    "wallet-top-up",
		Name:  "E-Wallet Top Up",
		Price: amount,
		Qty:   1,
	}}
	response, apiErr := c.snap.CreateTransaction(&snap.Request{
		TransactionDetails: midtransSDK.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: amount,
		},
		Items: &items,
	})
	if apiErr != nil {
		return "", "", fmt.Errorf("create Midtrans Snap transaction: %w", apiErr)
	}
	if response == nil || response.Token == "" || response.RedirectURL == "" {
		return "", "", fmt.Errorf("Midtrans returned an incomplete Snap transaction response")
	}

	return response.Token, response.RedirectURL, nil
}

func (c *Client) VerifyNotification(orderID, statusCode, grossAmount, signature string) bool {
	if c.serverKey == "" {
		return false
	}

	sum := sha512.Sum512([]byte(orderID + statusCode + grossAmount + c.serverKey))
	expected := []byte(hex.EncodeToString(sum[:]))
	provided := []byte(strings.ToLower(strings.TrimSpace(signature)))
	return len(expected) == len(provided) && subtle.ConstantTimeCompare(expected, provided) == 1
}

func (c *Client) GetTransactionStatus(orderID string) (*entity.PaymentStatus, error) {
	if c.serverKey == "" {
		return nil, fmt.Errorf("MIDTRANS_SERVER_KEY is not configured")
	}

	response, apiErr := c.core.CheckTransaction(orderID)
	if apiErr != nil {
		return nil, fmt.Errorf("check Midtrans transaction status: %w", apiErr)
	}
	if response == nil {
		return nil, fmt.Errorf("Midtrans returned an empty transaction status")
	}

	return &entity.PaymentStatus{
		OrderID:           response.OrderID,
		GrossAmount:       response.GrossAmount,
		Currency:          response.Currency,
		TransactionStatus: response.TransactionStatus,
		FraudStatus:       response.FraudStatus,
		PaymentType:       response.PaymentType,
	}, nil
}
