package http

import (
	accountHandler "transaction/internal/http/handler/account"
	userHandler "transaction/internal/http/handler/user"
	"transaction/internal/user/application"

	"github.com/labstack/echo/v4"
)

type Router struct {
	userHandler    *userHandler.Handler
	accountHandler *accountHandler.Handler
	userService    *application.Service
}

func NewRouter(userHandler *userHandler.Handler, accountHandler *accountHandler.Handler, userService *application.Service) *Router {
	return &Router{
		userHandler:    userHandler,
		accountHandler: accountHandler,
		userService:    userService,
	}
}

func (r *Router) Register(e *echo.Echo) {
	api := e.Group("/api/v1")

	api.POST("/users", r.userHandler.CreateUser)

	authAPI := api.Group("")
	authAPI.Use(AuthMiddleware(r.userService))

	authAPI.GET("/users/:id", r.userHandler.GetUser)
	authAPI.POST("/accounts", r.accountHandler.CreateAccount)
	authAPI.GET("/accounts", r.accountHandler.GetAccounts)
	authAPI.GET("/accounts/:id/balance", r.accountHandler.GetAccountBalance)
	authAPI.POST("/accounts/:id/deposit", r.accountHandler.Deposit)
	authAPI.POST("/transfers", r.accountHandler.Transfer)
	authAPI.GET("/accounts/:id/transactions", r.accountHandler.GetAccountTransactionHistory)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "healthy"})
	})
}
