package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/envm-vault/envm-vault/api/internal/config"
	"github.com/envm-vault/envm-vault/api/internal/database"
	"github.com/envm-vault/envm-vault/api/internal/email"
	authhandler "github.com/envm-vault/envm-vault/api/internal/handler/auth"
	healthhandler "github.com/envm-vault/envm-vault/api/internal/handler/health"
	userhandler "github.com/envm-vault/envm-vault/api/internal/handler/user"
	"github.com/envm-vault/envm-vault/api/internal/middleware"
	"github.com/envm-vault/envm-vault/api/internal/repository"
	"github.com/envm-vault/envm-vault/api/internal/router"
	"github.com/envm-vault/envm-vault/api/internal/service"
	authservice "github.com/envm-vault/envm-vault/api/internal/service/auth"
	userservice "github.com/envm-vault/envm-vault/api/internal/service/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	gormDB, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("get sql db: %v", err)
	}
	defer sqlDB.Close()

	userRepo := repository.NewUserRepository(gormDB)
	loginCodeRepo := repository.NewLoginCodeRepository(gormDB)

	var sender email.Sender
	if cfg.IsDevelopment() {
		sender = email.NewLoggerSender()
	} else {
		sender = email.NewNoopSender()
	}

	otpService := service.NewOTPService()
	authSvc := authservice.New(gormDB, userRepo, loginCodeRepo, otpService, sender, cfg.JWTSecret, cfg.JWTExpiry)
	userSvc := userservice.New(userRepo)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		},
	})

	app.Use(recover.New())
	app.Use(logger.New())

	router.Register(app, router.Deps{
		Health: healthhandler.New(gormDB),
		Auth:   authhandler.New(authSvc),
		User:   userhandler.New(userSvc),
		JWT:    middleware.NewJWTMiddleware(cfg.JWTSecret),
	})

	go func() {
		addr := ":" + cfg.Port
		log.Printf("server starting on %s (env=%s)", addr, cfg.AppEnv)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
}
