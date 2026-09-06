package audit

const (
	ChannelLoginAttempt  = "login_attempt"
	ChannelLoggedInUser  = "logged_in_user"
	ChannelPersonalToken = "personal_token"

	ActionLoginRequested            = "login.requested"
	ActionLoginSucceeded            = "login.succeeded"
	ActionLoginFailed               = "login.failed"
	ActionUserBootstrapSuperAdmin   = "user.bootstrap_super_admin"
	ActionUserRegistered            = "user.registered"
	ActionPersonalTokenCreated      = "personal_token.created"
	ActionPersonalTokenRotated      = "personal_token.rotated"
	ActionPersonalTokenRevoked      = "personal_token.revoked"
	ActionPersonalTokenRevokedAdmin = "personal_token.revoked_by_admin"

	SummaryLoginRequested            = "Login code requested"
	SummaryLoginSucceeded            = "Login succeeded"
	SummaryLoginFailed               = "Login failed"
	SummaryUserBootstrapSuperAdmin   = "Super admin bootstrapped"
	SummaryUserRegistered            = "User registered"
	SummaryPersonalTokenCreated      = "Personal token created"
	SummaryPersonalTokenRotated      = "Personal token rotated"
	SummaryPersonalTokenRevoked      = "Personal token revoked"
	SummaryPersonalTokenRevokedAdmin = "Personal token revoked by admin"
)
