package router

import (
	authhandler "github.com/envm-vault/envm-vault/api/internal/handler/auth"
	healthhandler "github.com/envm-vault/envm-vault/api/internal/handler/health"
	tokenhandler "github.com/envm-vault/envm-vault/api/internal/handler/token"
	userhandler "github.com/envm-vault/envm-vault/api/internal/handler/user"
	"github.com/envm-vault/envm-vault/api/internal/middleware"
)

type Deps struct {
	Health *healthhandler.Handler
	Auth   *authhandler.Handler
	User   *userhandler.Handler
	Token  *tokenhandler.Handler
	JWT    *middleware.JWTMiddleware
}
