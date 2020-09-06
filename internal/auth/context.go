package auth

import "context"

type currentUserContextKeyType struct{}

var currentUserContextKey currentUserContextKeyType


func ContextWithCurrentUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, currentUserContextKey, user)
}

func CurrentUserFromContext(ctx context.Context) *User {
	return ctx.Value(currentUserContextKey).(*User)
}
