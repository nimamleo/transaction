package http

import (
	"context"
	"fmt"
	"time"

	"transaction/pkg/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo   *echo.Echo
	router *Router
	config config.ServerConfig
}

func NewServer(cfg config.ServerConfig, router *Router) *Server {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	router.Register(e)

	return &Server{
		echo:   e,
		router: router,
		config: cfg,
	}
}

func (s *Server) Start() error {
	return s.echo.Start(fmt.Sprintf(":%s", s.config.Port))
}

func (s *Server) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return s.echo.Shutdown(shutdownCtx)
}
