package router

import (
	"net/http"

	authhandler "e-wallet/internal/handler/auth"
	userhandler "e-wallet/internal/handler/user"
	wallethandler "e-wallet/internal/handler/wallet"

	"github.com/gin-gonic/gin"
)

func RegisterHealthRoute(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

func RegisterUserRoutes(r *gin.Engine, handler *userhandler.Handler) {
	group := r.Group("/users")
	group.POST("/", handler.Create)
	group.GET("/", handler.GetAll)
	group.GET("/:id", handler.GetByID)
	group.PATCH("/:id/hide", handler.HideByID)
}

func RegisterAuthRoutes(r *gin.Engine, handler *authhandler.Handler) {
	r.POST("/auth/login", handler.Login)
}

func RegisterWalletRoutes(r *gin.Engine, handler *wallethandler.Handler) {
	group := r.Group("/wallets")
	group.GET("/user/:user_id", handler.GetByUserID)
	group.POST("/user/:user_id/topups", handler.CreateTopUp)
	r.POST("/payments/midtrans/notification", handler.MidtransNotification)
}
