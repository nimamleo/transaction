package http

import (
	userHandler "transaction/internal/http/handler/user"

	"github.com/labstack/echo/v4"
)

type Router struct {
	userHandler *userHandler.Handler
}

func NewRouter(userHandler *userHandler.Handler) *Router {
	return &Router{
		userHandler: userHandler,
	}
}

func (r *Router) Register(e *echo.Echo) {
	api := e.Group("/docs/v1")

	api.POST("/users", r.userHandler.CreateUser)
	api.GET("/users/:id", r.userHandler.GetUser)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "healthy"})
	})
}
