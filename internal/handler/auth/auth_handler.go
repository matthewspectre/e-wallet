package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"

	usecase "e-wallet/internal/usecase/auth"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc usecase.Usecase
}

func NewHandler(uc usecase.Usecase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}
	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}

	result, err := h.uc.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "email or password is incorrect"})
			return
		}
		log.Printf("failed to log in user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not log in"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		User: LoginUserDTO{
			ID:    result.User.ID,
			Name:  result.User.Name,
			Email: result.User.Email,
		},
	})
}
