package auth

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/envm-vault/envm-vault/api/internal/database"
	"github.com/envm-vault/envm-vault/api/internal/email"
	"github.com/envm-vault/envm-vault/api/internal/model"
	"github.com/envm-vault/envm-vault/api/internal/repository"
	"github.com/envm-vault/envm-vault/api/internal/service"
	userservice "github.com/envm-vault/envm-vault/api/internal/service/user"
)

type captureSender struct {
	lastEmail string
	lastCode  string
}

func (s *captureSender) SendLoginCode(_ context.Context, to, code string) error {
	s.lastEmail = to
	s.lastCode = code
	return nil
}

func TestService_BootstrapFlow(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	gormDB, err := database.Open(databaseURL)
	if err != nil {
		t.Fatalf("connect database: %v", err)
	}
	defer func() {
		sqlDB, _ := gormDB.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()

	if err := database.Ping(gormDB); err != nil {
		t.Skipf("database not available: %v", err)
	}

	ctx := context.Background()
	userRepo := repository.NewUserRepository(gormDB)
	loginCodeRepo := repository.NewLoginCodeRepository(gormDB)
	sender := &captureSender{}
	auth := New(gormDB, userRepo, loginCodeRepo, service.NewOTPService(), sender, "test-secret-key-for-jwt-signing", time.Hour)

	testEmail := "bootstrap-test@example.com"

	gormDB.WithContext(ctx).Where("email = ?", testEmail).Delete(&model.LoginCode{})
	gormDB.WithContext(ctx).Where("email = ?", testEmail).Delete(&model.User{})

	if err := auth.RequestLogin(ctx, testEmail); err != nil {
		t.Fatalf("RequestLogin() error = %v", err)
	}

	if sender.lastCode == "" {
		t.Fatal("expected OTP code to be sent")
	}

	resp, err := auth.VerifyLogin(ctx, testEmail, sender.lastCode)
	if err != nil {
		t.Fatalf("VerifyLogin() error = %v", err)
	}

	if resp.Token == "" {
		t.Error("expected JWT token")
	}
	if resp.User.Role != service.RoleSuperAdmin {
		t.Errorf("User.Role = %q, want %q", resp.User.Role, service.RoleSuperAdmin)
	}
	if resp.User.Email != testEmail {
		t.Errorf("User.Email = %q, want %q", resp.User.Email, testEmail)
	}

	gormDB.WithContext(ctx).Where("email = ?", testEmail).Delete(&model.LoginCode{})
	gormDB.WithContext(ctx).Where("email = ?", testEmail).Delete(&model.User{})
}

func TestService_RegisterUser(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	gormDB, err := database.Open(databaseURL)
	if err != nil {
		t.Fatalf("connect database: %v", err)
	}
	defer func() {
		sqlDB, _ := gormDB.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()

	if err := database.Ping(gormDB); err != nil {
		t.Skipf("database not available: %v", err)
	}

	ctx := context.Background()
	userRepo := repository.NewUserRepository(gormDB)
	loginCodeRepo := repository.NewLoginCodeRepository(gormDB)
	sender := &captureSender{}
	auth := New(gormDB, userRepo, loginCodeRepo, service.NewOTPService(), sender, "test-secret-key-for-jwt-signing", time.Hour)
	users := userservice.New(userRepo)

	adminEmail := "admin-register-test@example.com"
	newUserEmail := "new-user-register-test@example.com"

	gormDB.WithContext(ctx).Where("email IN ?", []string{adminEmail, newUserEmail}).Delete(&model.LoginCode{})
	gormDB.WithContext(ctx).Where("email IN ?", []string{adminEmail, newUserEmail}).Delete(&model.User{})

	if err := auth.RequestLogin(ctx, adminEmail); err != nil {
		t.Fatalf("RequestLogin() error = %v", err)
	}
	if _, err := auth.VerifyLogin(ctx, adminEmail, sender.lastCode); err != nil {
		t.Fatalf("VerifyLogin() error = %v", err)
	}

	user, err := users.Register(ctx, service.RoleSuperAdmin, newUserEmail, service.RoleUser)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if user.Email != newUserEmail {
		t.Errorf("user.Email = %q, want %q", user.Email, newUserEmail)
	}
	if user.Role != service.RoleUser {
		t.Errorf("user.Role = %q, want %q", user.Role, service.RoleUser)
	}

	_, err = users.Register(ctx, service.RoleUser, "another@example.com", service.RoleUser)
	if err != service.ErrForbidden {
		t.Errorf("Register() as user error = %v, want ErrForbidden", err)
	}

	gormDB.WithContext(ctx).Where("email IN ?", []string{adminEmail, newUserEmail}).Delete(&model.LoginCode{})
	gormDB.WithContext(ctx).Where("email IN ?", []string{adminEmail, newUserEmail}).Delete(&model.User{})
}

var _ email.Sender = (*captureSender)(nil)
