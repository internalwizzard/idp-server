package httpserver

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo *echo.Echo
}

func New() *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Recover())

	return &Server{echo: e}
}

func (s *Server) Echo() *echo.Echo {
	return s.echo
}

func (s *Server) Start(addr string) error {
	return s.echo.Start(addr)
}
