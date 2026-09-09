package router

import (
	audithandler "github.com/vmarinho22/envm-vault/api/internal/handler/audit"
	authhandler "github.com/vmarinho22/envm-vault/api/internal/handler/auth"
	environmenthandler "github.com/vmarinho22/envm-vault/api/internal/handler/environment"
	healthhandler "github.com/vmarinho22/envm-vault/api/internal/handler/health"
	projecthandler "github.com/vmarinho22/envm-vault/api/internal/handler/project"
	secrethandler "github.com/vmarinho22/envm-vault/api/internal/handler/secret"
	taghandler "github.com/vmarinho22/envm-vault/api/internal/handler/tag"
	teamhandler "github.com/vmarinho22/envm-vault/api/internal/handler/team"
	tokenhandler "github.com/vmarinho22/envm-vault/api/internal/handler/token"
	userhandler "github.com/vmarinho22/envm-vault/api/internal/handler/user"
	"github.com/vmarinho22/envm-vault/api/internal/middleware"
)

type Deps struct {
	Health      *healthhandler.Handler
	Auth        *authhandler.Handler
	User        *userhandler.Handler
	Token       *tokenhandler.Handler
	Audit       *audithandler.Handler
	Team        *teamhandler.Handler
	Tag         *taghandler.Handler
	Environment *environmenthandler.Handler
	Project     *projecthandler.Handler
	Secret      *secrethandler.Handler
	JWT         *middleware.JWTMiddleware
	APIKey      *middleware.APIKeyMiddleware
}
