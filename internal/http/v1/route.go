package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/magmaheat/merchStore/internal/service"
	"os"

	log "github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func New(services *service.Service) *echo.Echo {
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}", "method":"${method}","uri":"${uri}", "status":${status},"error":"${error}"}` + "\n",
		Output: setLogsFile(),
	}))
	e.Use(middleware.Recover())

	e.GET("/health", func(c echo.Context) error { return c.NoContent(200) })
	e.GET("/swagger/", echoSwagger.WrapHandler)

	authRoutes := NewAuthRoutes(services.Auth)
	e.POST("/api/auth", authRoutes.auth)

	authMiddleware := &AuthMiddleware{services.Auth}

	storeRoutes := NewStoreRoutes(services.Store)
	e.GET("/api/info", storeRoutes.info, authMiddleware.UserIdentity)
	e.POST("/api/sendCoin", storeRoutes.sendCoin, authMiddleware.UserIdentity)
	e.GET("/api/buy/:item", storeRoutes.buyItem, authMiddleware.UserIdentity)

	return e
}

func setLogsFile() *os.File {
	file, err := os.OpenFile("/logs/requests.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return file
}
