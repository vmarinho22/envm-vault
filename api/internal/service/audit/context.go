package audit

import "context"

type originContextKey struct{}

// RequestOrigin is network/auth channel metadata for the current request.
type RequestOrigin struct {
	IP          string
	UserAgent   string
	Origin      string
	AuthChannel string
}

func OriginFromContext(ctx context.Context) RequestOrigin {
	if ctx == nil {
		return RequestOrigin{AuthChannel: ChannelLoginAttempt}
	}
	origin, ok := ctx.Value(originContextKey{}).(RequestOrigin)
	if !ok {
		return RequestOrigin{AuthChannel: ChannelLoginAttempt}
	}
	if origin.AuthChannel == "" {
		origin.AuthChannel = ChannelLoginAttempt
	}
	return origin
}

func withOrigin(ctx context.Context, origin RequestOrigin) context.Context {
	if origin.AuthChannel == "" {
		origin.AuthChannel = ChannelLoginAttempt
	}
	return context.WithValue(ctx, originContextKey{}, origin)
}
