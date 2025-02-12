package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/magmaheat/merchStore/internal/service"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

const (
	userIdCtx = "userId"
)

type AuthMiddleware struct {
	authService service.Auth
}

func (a *AuthMiddleware) UserIdentity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := bearerToken(c.Request())
		if !ok {
			log.Errorf("http.AuthMiddleware.UserIdentity.bearerToken: %v", ErrInvalidAuthHeader)
			return newErrorResponse(c, http.StatusUnauthorized, ErrInvalidAuthHeader.Error())

		}

		userId, err := a.authService.ParseToken(token)
		if err != nil {
			log.Errorf("http.AuthMiddleware.UserIdentity.h.authService.ParseToken: %v", err)
			return newErrorResponse(c, http.StatusUnauthorized, ErrCannotParseToken.Error())
		}

		c.Set(userIdCtx, userId)

		return next(c)
	}
}

func bearerToken(r *http.Request) (string, bool) {
	const prefix = "Bearer "

	header := r.Header.Get(echo.HeaderAuthorization)
	if header == "" {
		return "", false
	}

	if len(header) > len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
		return header[len(prefix):], true
	}

	return "", false
}
