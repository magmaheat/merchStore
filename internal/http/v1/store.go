package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/magmaheat/merchStore/internal/service"
	"net/http"
)

type StoreRoutes struct {
	storeService service.Store
}

func NewStoreRoutes(g *echo.Group, storeService service.Store) {
	r := &StoreRoutes{
		storeService: storeService,
	}

	g.GET("/info", r.info)
	g.POST("/sendCoin", r.sendCoin)
	g.GET("/buy/:item", r.buyItem)
}

func (r *StoreRoutes) sendCoin(c echo.Context) error {
	return nil
}

func (r *StoreRoutes) info(c echo.Context) error {
	return nil
}

func (r *StoreRoutes) buyItem(c echo.Context) error {
	item := c.Param("item")
	if item == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "empty item")
	}

	return nil
}
