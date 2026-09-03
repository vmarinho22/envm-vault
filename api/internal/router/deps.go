package router

import (
	authhandler "github.com/envm-vault/envm-vault/api/internal/handler/auth"
	healthhandler "github.com/envm-vault/envm-vault/api/internal/handler/health"
	userhandler "github.com/envm-vault/envm-vault/api/internal/handler/user"
	"github.com/envm-vault/envm-vault/api/internal/middleware"
)

type Deps struct {
	Health *healthhandler.Handler
	Auth   *authhandler.Handler
	User   *userhandler.Handler
	JWT    *middleware.JWTMiddleware
}
