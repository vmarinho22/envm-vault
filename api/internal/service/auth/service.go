package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/envm-vault/envm-vault/api/internal/database"
	"github.com/envm-vault/envm-vault/api/internal/model"
	"github.com/envm-vault/envm-vault/api/internal/email"
	"github.com/envm-vault/envm-vault/api/internal/repository"
	"github.com/envm-vault/envm-vault/api/internal/service"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type Service struct {
	db           *gorm.DB
	users        *repository.UserRepository
	loginCodes   *repository.LoginCodeRepository
	otp          *service.OTPService
	email        email.Sender
	jwtSecret    []byte
	jwtExpiry    time.Duration
}

func New(
	db *gorm.DB,
	users *repository.UserRepository,
	loginCodes *repository.LoginCodeRepository,
	otp *service.OTPService,
	sender email.Sender,
	jwtSecret string,
	jwtExpiry time.Duration,
) *Service {
	return &Service{
		db:         db,
		users:      users,
		loginCodes: loginCodes,
		otp:        otp,
		email:      sender,
		jwtSecret:  []byte(jwtSecret),
		jwtExpiry:  jwtExpiry,
	}
}

func (s *Service) RequestLogin(ctx context.Context, rawEmail string) error {
	emailAddr, err := service.NormalizeEmail(rawEmail)
	if err != nil {
		return service.ErrInvalidEmail
	}

	count, err := s.users.Count(ctx)
	if err != nil {
		return err
	}

	if count > 0 {
		exists, err := s.users.ExistsByEmail(ctx, emailAddr)
		if err != nil {
			return err
		}
		if !exists {
			return service.ErrUserNotFound
		}
	}

	code, err := s.otp.Generate()
	if err != nil {
		return fmt.Errorf("generate otp: %w", err)
	}

	err = database.WithTx(ctx, s.db, func(tx *gorm.DB) error {
		loginCodes := s.loginCodes.WithDB(tx)
		if err := loginCodes.InvalidatePrevious(ctx, emailAddr); err != nil {
			return err
		}

		_, err := loginCodes.Create(ctx, &model.LoginCode{
			Email:     emailAddr,
			CodeHash:  s.otp.Hash(code),
			ExpiresAt: time.Now().Add(time.Duration(s.otp.ExpiryMinutes()) * time.Minute),
		})
		return err
	})
	if err != nil {
		return err
	}

	if err := s.email.SendLoginCode(ctx, emailAddr, code); err != nil {
		return fmt.Errorf("send login code: %w", err)
	}

	return nil
}

func (s *Service) VerifyLogin(ctx context.Context, rawEmail, code string) (*service.VerifyResponse, error) {
	emailAddr, err := service.NormalizeEmail(rawEmail)
	if err != nil {
		return nil, service.ErrInvalidEmail
	}

	if len(code) != 8 {
		return nil, service.ErrInvalidCode
	}

	loginCode, err := s.loginCodes.FindLatestValid(ctx, emailAddr)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrInvalidCode
		}
		return nil, err
	}

	if !s.otp.Verify(code, loginCode.CodeHash) {
		return nil, service.ErrInvalidCode
	}

	var user *model.User
	err = database.WithTx(ctx, s.db, func(tx *gorm.DB) error {
		loginCodes := s.loginCodes.WithDB(tx)
		users := s.users.WithDB(tx)

		if err := loginCodes.MarkUsed(ctx, loginCode.ID); err != nil {
			return err
		}

		found, err := users.FindByEmail(ctx, emailAddr)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if found != nil {
			user = found
			return nil
		}

		count, err := users.Count(ctx)
		if err != nil {
			return err
		}
		if count > 0 {
			return service.ErrUserNotFound
		}

		user, err = users.Create(ctx, &model.User{
			Email: emailAddr,
			Role:  service.RoleSuperAdmin,
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return nil, fmt.Errorf("generate jwt: %w", err)
	}

	return &service.VerifyResponse{
		Token: token,
		User:  service.ToUserResponse(*user),
	}, nil
}

func (s *Service) generateJWT(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(s.jwtExpiry).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
