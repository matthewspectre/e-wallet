package user

import (
	"net/http"
	"strings"
	"time"

	entity "e-wallet/internal/entity/user"
	usecase "e-wallet/internal/usecase/user"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	uc usecase.Usecase
}

func NewHandler(uc usecase.Usecase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) Create(c *gin.Context) {
	var req UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.TrimSpace(req.Name) == "" ||
		strings.TrimSpace(req.Email) == "" ||
		req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, email, and password are required"})
		return
	}
	if len(req.Password) > 72 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must not exceed 72 bytes"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	now := time.Now()
	data := &entity.User{
		Name:         strings.TrimSpace(req.Name),
		Email:        strings.TrimSpace(req.Email),
		PasswordHash: string(passwordHash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := h.uc.Create(data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, UserResponse{
		ID:        data.ID,
		Name:      data.Name,
		Email:     data.Email,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	})
}
