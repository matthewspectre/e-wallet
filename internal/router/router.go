package router

import (
	"net/http"

	userhandler "e-wallet/internal/handler/user"

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
}
