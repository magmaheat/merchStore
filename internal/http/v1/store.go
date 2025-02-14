package v1

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/magmaheat/merchStore/internal/service"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type StoreRoutes struct {
	storeService service.Store
}

func NewStoreRoutes(storeService service.Store) *StoreRoutes {
	return &StoreRoutes{storeService: storeService}
}

type inputSendCoin struct {
	ToUser string `json:"toUser" validate:"required"`
	Amount int    `json:"amount" validate:"required"`
}

func (r *StoreRoutes) sendCoin(c echo.Context) error {
	var input inputSendCoin

	if err := c.Bind(&input); err != nil {
		log.Errorf("http.v1.sendCoin.Bind: %v", err)
		return newErrorResponse(c, http.StatusBadRequest, ErrInvalidRequestBody.Error())
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		log.Errorf("http.v1.sendCoin.Struct: %v", err)
		return newErrorResponse(c, http.StatusBadRequest, ErrInvalidRequestBody.Error())
	}

	err := r.storeService.SendCoin(c.Request().Context(), input.ToUser, input.Amount)
	if err != nil {
		return newErrorResponse(c, http.StatusInternalServerError, ErrInternalServerError.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (r *StoreRoutes) info(c echo.Context) error {
	return nil
}

func (r *StoreRoutes) buyItem(c echo.Context) error {
	item := c.Param("item")
	if item == "" {
		return newErrorResponse(c, http.StatusBadRequest, "empty item")
	}

	err := r.storeService.BuyItem(c.Request().Context(), item)
	if err != nil {
		if errors.Is(err, service.ErrItemNotFound) {
			return newErrorResponse(c, http.StatusBadRequest, err.Error())
		}

		if errors.Is(err, service.ErrBalanceTooLow) {
			return newErrorResponse(c, http.StatusBadRequest, err.Error())
		}

		return newErrorResponse(c, http.StatusInternalServerError, ErrInternalServerError.Error())
	}

	return c.NoContent(http.StatusOK)
}
