package router

import (
	audithandler "github.com/envm-vault/envm-vault/api/internal/handler/audit"
	authhandler "github.com/envm-vault/envm-vault/api/internal/handler/auth"
	healthhandler "github.com/envm-vault/envm-vault/api/internal/handler/health"
	taghandler "github.com/envm-vault/envm-vault/api/internal/handler/tag"
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
	Tag    *taghandler.Handler
	JWT    *middleware.JWTMiddleware
}
