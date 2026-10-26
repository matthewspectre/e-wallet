package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	registrationEmail "e-wallet/internal/email"
	handlerAuth "e-wallet/internal/handler/auth"
	handlerUser "e-wallet/internal/handler/user"
	handlerWallet "e-wallet/internal/handler/wallet"
	"e-wallet/internal/midtrans"
	mysqlUser "e-wallet/internal/mysql/user"
	mysqlWallet "e-wallet/internal/mysql/wallet"
	"e-wallet/internal/router"
	usecaseAuth "e-wallet/internal/usecase/auth"
	usecaseUser "e-wallet/internal/usecase/user"
	usecaseWallet "e-wallet/internal/usecase/wallet"
)

func main() {
	var httpPortFlag string
	flag.StringVar(&httpPortFlag, "http-port", "", "HTTP listen port (overrides APP_PORT/PORT env var)")
	flag.Parse()

	dbUser := strings.TrimSpace(os.Getenv("DB_USER"))
	dbPass := os.Getenv("DB_PASS")
	dbHost := strings.TrimSpace(os.Getenv("DB_HOST"))
	dbPort := strings.TrimSpace(os.Getenv("DB_PORT"))
	dbName := strings.TrimSpace(os.Getenv("DB_NAME"))

	if dbUser == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("database environment variables (DB_USER, DB_HOST, DB_PORT, DB_NAME) must not be empty")
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to initialize database connection: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	defer sqlDB.Close()

	userRepo := mysqlUser.NewRepository(db)
	userUC := usecaseUser.NewUsecase(userRepo)
	authUC := usecaseAuth.NewUsecase(userRepo)
	authHTTPHandler := handlerAuth.NewHandler(authUC)
	emailFrom := strings.TrimSpace(os.Getenv("EMAIL_FROM"))
	if emailFrom == "" {
		emailFrom = "matthewxrplb@gmail.com"
	}
	emailSender := registrationEmail.NewSender(emailFrom, os.Getenv("EMAIL_APP_PASSWORD"))
	userHTTPHandler := handlerUser.NewHandler(userUC, emailSender)

	walletRepo := mysqlWallet.NewRepository(db)
	midtransClient := midtrans.NewClient(os.Getenv("MIDTRANS_SERVER_KEY"))
	walletUC := usecaseWallet.NewUsecase(walletRepo, midtransClient)
	walletHTTPHandler := handlerWallet.NewHandler(walletUC)

	r := gin.Default()
	allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
	if strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS")) == "" {
		allowedOrigins = []string{"http://localhost:5173", "http://127.0.0.1:5173"}
	}
	for i := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{"GET", "POST", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type"},
		MaxAge:       12 * time.Hour,
	}))
	router.RegisterHealthRoute(r)
	router.RegisterAuthRoutes(r, authHTTPHandler)
	router.RegisterUserRoutes(r, userHTTPHandler)
	router.RegisterWalletRoutes(r, walletHTTPHandler)

	httpPort := strings.TrimSpace(httpPortFlag)
	if httpPort == "" {
		httpPort = strings.TrimSpace(os.Getenv("APP_PORT"))
	}
	if httpPort == "" {
		httpPort = strings.TrimSpace(os.Getenv("PORT"))
	}
	if httpPort == "" {
		httpPort = "8080"
	}

	listenAddr := httpPort
	if !strings.Contains(listenAddr, ":") {
		listenAddr = ":" + listenAddr
	}

	log.Printf("listening and serving HTTP on %s", listenAddr)
	if err := r.Run(listenAddr); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
