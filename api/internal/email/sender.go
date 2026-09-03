package email

import "context"

type Sender interface {
	SendLoginCode(ctx context.Context, to, code string) error
}
