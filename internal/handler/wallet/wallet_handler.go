package wallet

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	entity "e-wallet/internal/entity/wallet"
	usecase "e-wallet/internal/usecase/wallet"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	uc usecase.Usecase
}

func NewHandler(uc usecase.Usecase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) GetByUserID(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id must be a positive integer"})
		return
	}

	data, err := h.uc.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "wallet not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toWalletResponse(data))
}

func (h *Handler) CreateTopUp(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id must be a positive integer"})
		return
	}

	var req TopUpCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	topUp, err := h.uc.CreateTopUp(userID, req.Amount, req.PaymentMethod)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "wallet not found"})
			return
		}
		if errors.Is(err, usecase.ErrInvalidTopUpAmount) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, usecase.ErrInvalidPaymentMethod) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("failed to create top-up for user %d: %v", userID, err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "could not create Midtrans payment"})
		return
	}

	c.JSON(http.StatusCreated, TopUpResponse{
		ID:          topUp.ID,
		OrderID:     topUp.OrderID,
		Amount:      topUp.Amount,
		Currency:    "IDR",
		Status:      topUp.Status,
		SnapToken:   topUp.SnapToken,
		RedirectURL: topUp.RedirectURL,
		CreatedAt:   topUp.CreatedAt,
	})
}

func (h *Handler) MidtransNotification(c *gin.Context) {
	var req MidtransNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.uc.HandleNotification(req.OrderID, req.StatusCode, req.GrossAmount, req.SignatureKey)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidNotification) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Midtrans notification signature"})
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "top-up transaction not found"})
			return
		}
		if errors.Is(err, usecase.ErrUnexpectedPayment) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Midtrans transaction does not match top-up"})
			return
		}
		log.Printf("failed to process Midtrans notification for order %s: %v", req.OrderID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not process Midtrans notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification processed"})
}

func toWalletResponse(data *entity.Wallet) WalletResponse {
	return WalletResponse{
		ID:        data.ID,
		UserID:    data.UserID,
		Balance:   data.Balance,
		Currency:  data.Currency,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}
