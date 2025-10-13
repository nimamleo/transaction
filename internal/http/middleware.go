package http

import (
	"transaction/internal/user/application"
	"transaction/pkg/httpcontext"
	"transaction/pkg/stdresponse"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(userService *application.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			apiKey := c.Request().Header.Get("X-API-KEY")
			if apiKey == "" {
				return stdresponse.SendHttpResponse(c, "missing API key")
			}

			userID, err := userService.GetUserIDByAPIKey(c.Request().Context(), apiKey)
			if err != nil {
				return stdresponse.SendHttpResponse(c, err)
			}

			user, err := userService.GetUserByID(c.Request().Context(), userID)
			if err != nil {
				return stdresponse.SendHttpResponse(c, err)
			}

			ctx := httpcontext.SetUser(c.Request().Context(), user)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
