package router

import (
	authhandler "github.com/envm-vault/envm-vault/api/internal/handler/auth"
	audithandler "github.com/envm-vault/envm-vault/api/internal/handler/audit"
	healthhandler "github.com/envm-vault/envm-vault/api/internal/handler/health"
	teamhandler "github.com/envm-vault/envm-vault/api/internal/handler/team"
	tokenhandler "github.com/envm-vault/envm-vault/api/internal/handler/token"
	userhandler "github.com/envm-vault/envm-vault/api/internal/handler/user"
	"github.com/envm-vault/envm-vault/api/internal/middleware"
)

type Deps struct {
	Health *healthhandler.Handler
	Auth   *authhandler.Handler
	User   *userhandler.Handler
	Token  *tokenhandler.Handler
	Audit  *audithandler.Handler
	Team   *teamhandler.Handler
	JWT    *middleware.JWTMiddleware
}
