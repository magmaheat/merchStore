package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/magmaheat/merchStore/internal/service"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type AuthRoutes struct {
	authService service.Auth
}

func NewAuthRoutes(g *echo.Group, authService service.Auth) {
	r := &AuthRoutes{authService: authService}

	g.POST("/auth", r.auth)
}

type authRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (r *AuthRoutes) auth(c echo.Context) error {
	var input authRequest

	if err := c.Bind(&input); err != nil {
		log.Error("AuthRoutes.c.Bind: Failed to bind auth request")
		return newErrorResponse(c, http.StatusBadRequest, ErrInvalidRequestBody.Error())
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return newErrorResponse(c, http.StatusBadRequest, ErrInvalidRequestBody.Error())
	}

	token, err := r.authService.GenerateToken(c.Request().Context(), service.AuthGenerateTokenInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		return newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	type response struct {
		Token string `json:"token"`
	}

	return c.JSON(http.StatusOK, response{
		Token: token,
	})
}
