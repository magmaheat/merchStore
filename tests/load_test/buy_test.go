//go:build buy

package load_test

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"golang.org/x/time/rate"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestBuyLoad(t *testing.T) {
	tokens, err := registerUsersParallel(100000)
	if err != nil {
		t.Fatalf("Failed to register users: %v", err)
	}

	rand.Seed(time.Now().UnixNano())

	rate := vegeta.Rate{Freq: 50, Per: time.Second}
	duration := 10 * time.Second

	attacker := vegeta.NewAttacker()

	targeter := func() vegeta.Targeter {
		return func(tgt *vegeta.Target) error {
			if tgt == nil {
				return vegeta.ErrNilTarget
			}

			token := tokens[rand.Intn(len(tokens))]

			tgt.Method = "GET"
			tgt.URL = "http://localhost:8080/api/buy/pen"
			tgt.Header = map[string][]string{
				"Authorization": {fmt.Sprintf("Bearer %s", token)},
			}
			return nil
		}
	}

	var metrics vegeta.Metrics

	for res := range attacker.Attack(targeter(), rate, duration, "Buy Test") {
		metrics.Add(res)
	}

	metrics.Close()

	t.Logf("Requests: %d\n", metrics.Requests)
	t.Logf("Success: %.2f%%\n", metrics.Success*100)
	t.Logf("Latency (mean): %s\n", metrics.Latencies.Mean)
	t.Logf("Latency (99th percentile): %s\n", metrics.Latencies.P99)
	t.Logf("Status codes: %v\n", metrics.StatusCodes)

	// Проверка SLI успешности.
	if metrics.Success < 0.9999 {
		t.Error("SLI успешности не выполнен (меньше 99.99%)")
	} else {
		t.Log("SLI успешности выполнен")
	}

	if metrics.Latencies.Mean > 50*time.Millisecond {
		t.Error("❌ SLI времени ответа не выполнен (больше 50 мс)")
	} else {
		t.Log("✅ SLI времени ответа выполнен")
	}
}

func generateUserList(count int) []string {
	users := make([]string, 0, count)
	for i := 0; i < count; i++ {
		users = append(users, fmt.Sprintf("user_%d", i+1))
	}

	return users
}

func registerUsers(count int) ([]string, error) {
	users := generateUserList(count)
	tokens := make([]string, 0, count)

	for idx, user := range users {
		resp, err := http.Post(
			"http://localhost:8080/api/auth",
			"application/json",
			strings.NewReader(fmt.Sprintf(`{"username": "%s", "password": "test_pass"}`, user)),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to register user %s: %w", user, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to register user %s: status %d", user, resp.StatusCode)
		}

		var result struct {
			Token string `json:"token"`
		}
		if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode response for user %s: %w", user, err)
		}

		tokens = append(tokens, result.Token)

		if idx > 0 && idx%500 == 0 {
			log.Infof("registered: %d users", idx)
		}
	}

	return tokens, nil
}

func registerUsersParallel(count int) ([]string, error) {
	users := generateUserList(count)
	tokens := make([]string, count)
	var wg sync.WaitGroup
	errChan := make(chan error, count)

	// Запускаем горутину для чтения ошибок.
	var errors []error
	go func() {
		for err := range errChan {
			errors = append(errors, err)
		}
	}()

	// Ограничиваем количество запросов до 100 в секунду.
	limiter := rate.NewLimiter(100, 1)

	for i, user := range users {
		wg.Add(1)
		go func(i int, user string) {
			defer wg.Done()

			// Ожидаем разрешения от rate limiter.
			if err := limiter.Wait(context.Background()); err != nil {
				errChan <- fmt.Errorf("rate limiter error for user %s: %w", user, err)
				return
			}

			resp, err := http.Post(
				"http://localhost:8080/api/auth",
				"application/json",
				strings.NewReader(fmt.Sprintf(`{"username": "%s", "password": "test_pass"}`, user)),
			)
			if err != nil {
				errChan <- fmt.Errorf("failed to register user %s: %w", user, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body) // Читаем тело ответа для диагностики.
				errChan <- fmt.Errorf(
					"failed to register user %s: status %d, response: %s",
					user, resp.StatusCode, string(body),
				)
				return
			}

			var result struct {
				Token string `json:"token"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				errChan <- fmt.Errorf("failed to decode response for user %s: %w", user, err)
				return
			}

			tokens[i] = result.Token
		}(i, user)
	}

	wg.Wait()
	close(errChan) // Закрываем канал после завершения всех горутин.

	// Если есть ошибки, возвращаем первую из них.
	if len(errors) > 0 {
		return nil, errors[0]
	}

	log.Info("all users are registered")
	return tokens, nil
}
