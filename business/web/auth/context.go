package auth

import (
	"context"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// key is used to store/retrieve a Claims value from a context.Context.
const claimKey ctxKey = 1

// key is used to store/retrieve a user value from a context.Context.
const userKey ctxKey = 2

// =============================================================================

// SetClaims stores the claims in the context.
func SetClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, claimKey, claims)
}

// GetClaims returns the claims from the context.
func GetClaims(ctx context.Context) Claims {
	v, ok := ctx.Value(claimKey).(Claims)
	if !ok {
		return Claims{}
	}
	return v
}

// SetUserID stores the user id from the request in the context.
func SetUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, userKey, userID)
}

// GetUserID returns the claims from the context.
func GetUserID(ctx context.Context) int {
	v, ok := ctx.Value(userKey).(int)
	if !ok {
		return 0
	}
	return v
}
