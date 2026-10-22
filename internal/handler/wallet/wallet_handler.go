package wallet

import (
	"errors"
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
