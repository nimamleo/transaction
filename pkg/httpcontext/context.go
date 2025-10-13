package httpcontext

import (
	"context"

	"transaction/internal/user/domain"

	"github.com/labstack/echo/v4"
)

type contextKey string

const UserKey contextKey = "user"

func SetUser(ctx context.Context, user *domain.User) context.Context {
	return context.WithValue(ctx, UserKey, user)
}

func GetUser(c echo.Context) *domain.User {
	user, ok := c.Request().Context().Value(UserKey).(*domain.User)
	if !ok {
		return nil
	}
	return user
}
