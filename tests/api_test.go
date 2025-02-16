package tests

import (
	"github.com/gavv/httpexpect/v2"
	"net/http"
	"testing"
)

func TestBuyItem(t *testing.T) {
	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	var token string

	t.Run("Auth", func(t *testing.T) {
		authReq := map[string]string{
			"username": "testUser",
			"password": "testPassword",
		}

		obj := e.POST("/api/auth").
			WithJSON(authReq).
			Expect().
			Status(http.StatusOK).JSON().Object()

		token = obj.Value("token").String().Raw()
		if token == "" {
			t.Fatal("Expected non-empty token")
		}
	})

	t.Run("BuyItemSuccess", func(t *testing.T) {
		e.GET("/api/buy/cup").WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK)
	})

	t.Run("BuyItemFail", func(t *testing.T) {
		e.GET("/api/buy/board").WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusBadRequest)
	})
}

func TestSendCoin(t *testing.T) {
	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	var token1 string

	t.Run("AuthUser1", func(t *testing.T) {
		authReq := map[string]string{
			"username": "testUser1",
			"password": "testPassword",
		}

		obj := e.POST("/api/auth").
			WithJSON(authReq).
			Expect().
			Status(http.StatusOK).JSON().Object()

		token1 = obj.Value("token").String().Raw()
		if token1 == "" {
			t.Fatal("Expected non-empty token")
		}
	})

	var token2 string

	t.Run("AuthUser2", func(t *testing.T) {
		authReq := map[string]string{
			"username": "testUser2",
			"password": "testPassword",
		}

		obj := e.POST("/api/auth").
			WithJSON(authReq).
			Expect().
			Status(http.StatusOK).JSON().Object()

		token2 = obj.Value("token").String().Raw()
		if token2 == "" {
			t.Fatal("Expected non-empty token")
		}
	})

	t.Run("SendCoin", func(t *testing.T) {
		sendCoinReq := map[string]interface{}{
			"toUser": "testUser2",
			"amount": 10,
		}
		e.POST("/api/sendCoin").WithHeader("Authorization", "Bearer "+token1).WithJSON(sendCoinReq).
			Expect().
			Status(http.StatusOK)
	})

	t.Run("SendCoin", func(t *testing.T) {
		sendCoinReq := map[string]interface{}{
			"toUser": "testUser2",
			"amount": 1000,
		}
		e.POST("/api/sendCoin").WithHeader("Authorization", "Bearer "+token1).WithJSON(sendCoinReq).
			Expect().
			Status(http.StatusInternalServerError)
	})
}

func TestGetInfo(t *testing.T) {
	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	var token string

	t.Run("Auth", func(t *testing.T) {
		authReq := map[string]string{
			"username": "testUser",
			"password": "testPassword",
		}

		obj := e.POST("/api/auth").
			WithJSON(authReq).
			Expect().
			Status(http.StatusOK).JSON().Object()

		token = obj.Value("token").String().Raw()
		if token == "" {
			t.Fatal("Expected non-empty token")
		}
	})

	t.Run("GetInfoSuccess", func(t *testing.T) {
		e.GET("/api/info").WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).JSON().Object().ContainsKey("coins").
			Value("coins").Number().Gt(0)
	})

	t.Run("GetInfoUnauthorized", func(t *testing.T) {
		e.GET("/api/info").
			Expect().
			Status(http.StatusUnauthorized)
	})
}
