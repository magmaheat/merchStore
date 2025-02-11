package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/magmaheat/merchStore/internal/service"
)

func New(services *service.Service) *echo.Echo {
	e := echo.New()

	return e
}
