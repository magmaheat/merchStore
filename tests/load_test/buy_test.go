//go:build buy

package load_test

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"golang.org/x/exp/rand"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestBuyLoad(t *testing.T) {
	tokens, err := registerUsers(100000)
	if err != nil {
		t.Fatalf("Failed to register users: %v", err)
	}

	rand.Seed(uint64(time.Now().UnixNano()))

	rate := vegeta.Rate{Freq: 1000, Per: time.Second}
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

	if metrics.Success < 0.9999 {
		t.Error("❌ SLI успешности не выполнен (меньше 99.99%)")
	} else {
		t.Log("✅ SLI времени ответа выполнен")
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
	timeStart := time.Now()
	users := generateUserList(count)
	tokens := make([]string, 0, count)
	var mu sync.Mutex

	// количество горутин, отправляющих запрросы
	workerLimit := 100
	sem := make(chan struct{}, workerLimit)
	errCh := make(chan error, count)
	var wg sync.WaitGroup

	for idx, user := range users {
		sem <- struct{}{}
		wg.Add(1)

		go func(user string, idx int) {
			defer wg.Done()
			defer func() { <-sem }()

			token, err := registerUser(user)
			if err != nil {
				errCh <- err
				return
			}

			mu.Lock()
			tokens = append(tokens, token)
			mu.Unlock()

			if idx > 0 && idx%500 == 0 {
				log.Infof("registered: %d users", idx)
			}
		}(user, idx)
	}

	wg.Wait()
	close(errCh)

	if len(errCh) > 0 {
		return nil, <-errCh
	}

	log.Printf("Общее время регистрации пользователей: %v", time.Now().Sub(timeStart))

	return tokens, nil
}

func registerUser(user string) (string, error) {
	resp, err := http.Post(
		"http://localhost:8080/api/auth",
		"application/json",
		strings.NewReader(fmt.Sprintf(`{"username": "%s", "password": "test_pass"}`, user)),
	)
	if err != nil {
		return "", fmt.Errorf("failed to register user %s: %w", user, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to register user %s: status %d", user, resp.StatusCode)
	}

	var result struct {
		Token string `json:"token"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response for user %s: %w", user, err)
	}

	return result.Token, nil
}
