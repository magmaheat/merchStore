package v1

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/magmaheat/merchStore/internal/service"
	"net/http"
)

type StoreRoutes struct {
	storeService service.Store
}

func NewStoreRoutes(storeService service.Store) *StoreRoutes {
	return &StoreRoutes{storeService: storeService}
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

	err := r.storeService.BuyItem(c.Request().Context(), item)
	if err != nil {
		if errors.Is(err, service.ErrItemNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "item not found")
		}
		if errors.Is(err, service.ErrBalanceTooLow) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
