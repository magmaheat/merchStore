package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/magmaheat/merchStore/internal/service"
)

type StoreRoutes struct {
	storeService service.Store
}

func NewStoreRoutes(g *echo.Group, storeService service.Store) {
	r := &StoreRoutes{
		storeService: storeService,
	}

	_ = r
}
