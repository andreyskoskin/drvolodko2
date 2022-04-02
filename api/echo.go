package api

import (
	"github.com/labstack/echo/v4"
)

type (
	EchoConfig struct {
		Address string `toml:"address"`
	}

	EchoAPI interface {
		Bind(g *echo.Group)
	}

	EchoBindings map[string]EchoAPI

	Echo struct {
		e *echo.Echo
	}
)

func StartEcho(config EchoConfig, bindings EchoBindings) error {
	var e = echo.New()

	for prefix, api := range bindings {
		api.Bind(e.Group(prefix))
	}

	return e.Start(config.Address)
}
