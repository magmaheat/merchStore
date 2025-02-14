package v1

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

var (
	ErrInvalidAuthHeader   = fmt.Errorf("invalid auth header")
	ErrCannotParseToken    = fmt.Errorf("cannot parse token")
	ErrInvalidRequestBody  = fmt.Errorf("invalid request body")
	ErrInternalServerError = fmt.Errorf("internal server error")
)

func newErrorResponse(c echo.Context, code int, msg string) error {
	return c.JSON(code, map[string]string{"error": msg})
}
