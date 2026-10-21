package user

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	entity "e-wallet/internal/entity/user"
	usecase "e-wallet/internal/usecase/user"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
		NoTelp:       strings.TrimSpace(req.NoTelp),
		PasswordHash: string(passwordHash),
		Visible:      1,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := h.uc.Create(data); err != nil {
		if message, duplicate := duplicateMessage(err); duplicate {
			c.JSON(http.StatusConflict, gin.H{
				"status":  http.StatusConflict,
				"message": message,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, UserResponse{
		ID:        data.ID,
		Name:      data.Name,
		Email:     data.Email,
		NoTelp:    data.NoTelp,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	})
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive integer"})
		return
	}

	data, err := h.uc.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toUserResponse(data))
}

func (h *Handler) GetAll(c *gin.Context) {
	users, total, err := h.uc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]UserResponse, 0, len(users))
	for _, data := range users {
		responses = append(responses, toUserResponse(data))
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  responses,
		"total": total,
	})
}

func (h *Handler) HideByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive integer"})
		return
	}

	if err := h.uc.HideByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user hidden"})
}

func toUserResponse(data *entity.User) UserResponse {
	return UserResponse{
		ID:        data.ID,
		Name:      data.Name,
		Email:     data.Email,
		NoTelp:    data.NoTelp,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}

func duplicateMessage(err error) (string, bool) {
	errMessage := strings.ToLower(err.Error())
	if !strings.Contains(errMessage, "duplicate entry") &&
		!strings.Contains(errMessage, "error 1062") {
		return "", false
	}

	emailExists := strings.Contains(errMessage, "email")
	phoneExists := strings.Contains(errMessage, "no_telp") ||
		strings.Contains(errMessage, "phone")

	switch {
	case emailExists && phoneExists:
		return "email and phone number already exist", true
	case emailExists:
		return "email already exists", true
	case phoneExists:
		return "phone number already exists", true
	default:
		return "email or phone number already exists", true
	}
}
