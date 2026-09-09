package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/envm-vault/envm-vault/api/internal/database"
	"github.com/envm-vault/envm-vault/api/internal/email"
	"github.com/envm-vault/envm-vault/api/internal/env"
	audithandler "github.com/envm-vault/envm-vault/api/internal/handler/audit"
	authhandler "github.com/envm-vault/envm-vault/api/internal/handler/auth"
	environmenthandler "github.com/envm-vault/envm-vault/api/internal/handler/environment"
	healthhandler "github.com/envm-vault/envm-vault/api/internal/handler/health"
	"github.com/envm-vault/envm-vault/api/internal/handler/httperr"
	projecthandler "github.com/envm-vault/envm-vault/api/internal/handler/project"
	secrethandler "github.com/envm-vault/envm-vault/api/internal/handler/secret"
	taghandler "github.com/envm-vault/envm-vault/api/internal/handler/tag"
	teamhandler "github.com/envm-vault/envm-vault/api/internal/handler/team"
	tokenhandler "github.com/envm-vault/envm-vault/api/internal/handler/token"
	userhandler "github.com/envm-vault/envm-vault/api/internal/handler/user"
	"github.com/envm-vault/envm-vault/api/internal/middleware"
	"github.com/envm-vault/envm-vault/api/internal/repository"
	"github.com/envm-vault/envm-vault/api/internal/router"
	"github.com/envm-vault/envm-vault/api/internal/secretcrypto"
	"github.com/envm-vault/envm-vault/api/internal/service"
	auditservice "github.com/envm-vault/envm-vault/api/internal/service/audit"
	authservice "github.com/envm-vault/envm-vault/api/internal/service/auth"
	environmentservice "github.com/envm-vault/envm-vault/api/internal/service/environment"
	projectservice "github.com/envm-vault/envm-vault/api/internal/service/project"
	secretservice "github.com/envm-vault/envm-vault/api/internal/service/secret"
	tagservice "github.com/envm-vault/envm-vault/api/internal/service/tag"
	teamservice "github.com/envm-vault/envm-vault/api/internal/service/team"
	tokenservice "github.com/envm-vault/envm-vault/api/internal/service/token"
	userservice "github.com/envm-vault/envm-vault/api/internal/service/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	if err := env.Load(); err != nil {
		log.Fatalf("load env: %v", err)
	}

	crypter, err := secretcrypto.NewAESGCM(env.Envs.Crypto.Key, env.Envs.Crypto.Version)
	if err != nil {
		log.Fatalf("init encryption: %v", err)
	}

	gormDB, err := database.Open(env.Envs.Database.URL, crypter)
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
	tokenRepo := repository.NewAccessTokenRepository(gormDB)
	auditRepo := repository.NewAuditLogRepository(gormDB)
	teamRepo := repository.NewTeamRepository(gormDB)
	tagRepo := repository.NewTagRepository(gormDB)
	environmentRepo := repository.NewEnvironmentRepository(gormDB)
	projectRepo := repository.NewProjectRepository(gormDB)
	secretRepo := repository.NewSecretRepository(gormDB)
	auditSvc := auditservice.New(auditRepo, userRepo)

	var sender email.Sender
	if env.Envs.IsDevelopment() {
		sender = email.NewLoggerSender()
	} else {
		sender = email.NewNoopSender()
	}

	otpService := service.NewOTPService()
	authSvc := authservice.New(gormDB, userRepo, loginCodeRepo, otpService, sender, env.Envs.JWT.Secret, env.Envs.JWT.Expiry, auditSvc)
	userSvc := userservice.New(userRepo, tokenRepo, auditSvc)
	tokenSvc := tokenservice.New(tokenRepo, userRepo, auditSvc)
	teamSvc := teamservice.New(teamRepo, userRepo, auditSvc)
	tagSvc := tagservice.New(tagRepo, auditSvc)
	environmentSvc := environmentservice.New(environmentRepo, auditSvc)
	projectSvc := projectservice.New(projectRepo, teamRepo, environmentRepo, tagRepo, secretRepo, auditSvc)
	secretSvc := secretservice.New(secretRepo, projectRepo, crypter, auditSvc)

	app := fiber.New(fiber.Config{
		ErrorHandler: httperr.AppHandler,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(env.Envs.CORS.Origins, ","),
		AllowHeaders: "Origin, Content-Type, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))
	app.Use(middleware.NewAudit(auditSvc))
	app.Use(middleware.NewRateLimiter(env.Envs.RateLimit.Max, env.Envs.RateLimit.Window))

	router.Register(app, router.Deps{
		Health:      healthhandler.New(gormDB),
		Auth:        authhandler.New(authSvc),
		User:        userhandler.New(userSvc),
		Token:       tokenhandler.New(tokenSvc),
		Audit:       audithandler.New(auditSvc),
		Team:        teamhandler.New(teamSvc),
		Tag:         taghandler.New(tagSvc),
		Environment: environmenthandler.New(environmentSvc),
		Project:     projecthandler.New(projectSvc),
		Secret:      secrethandler.New(secretSvc),
		JWT:         middleware.NewJWTMiddleware(env.Envs.JWT.Secret, userRepo, auditSvc),
		APIKey:      middleware.NewAPIKeyMiddleware(tokenRepo, userRepo, auditSvc),
	})
	router.MountFrontend(app, env.Envs.App.FrontendDist)

	go func() {
		addr := ":" + env.Envs.App.Port
		log.Printf("server starting on %s (env=%s)", addr, env.Envs.App.Env)
		if dist := env.Envs.App.FrontendDist; dist != "" {
			log.Printf("serving frontend from %s", dist)
		}
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
