//go:build load

package load_test

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func generateUsername() string {
	return fmt.Sprintf("user_%s", uuid.New().String())
}

func generateUserList(count int) []string {
	users := make([]string, 0, count)
	for i := 0; i < count; i++ {
		users = append(users, fmt.Sprintf("user_%d", i+1))
	}

	return users
}

func TestAuthLoad(t *testing.T) {
	users := generateUserList(100000)
	var index uint64

	rate := vegeta.Rate{Freq: 50, Per: time.Second}
	duration := 10 * time.Second

	attacker := vegeta.NewAttacker()

	targeter := func() vegeta.Targeter {
		return func(tgt *vegeta.Target) error {
			if tgt == nil {
				return vegeta.ErrNilTarget
			}

			currentIndex := atomic.AddUint64(&index, 1) - 1
			user := users[int(currentIndex)%len(users)]

			tgt.Method = "POST"
			tgt.URL = "http://localhost:8080/api/auth"
			tgt.Body = []byte(fmt.Sprintf(`{"username": "%s", "password": "test_pass"}`, user))
			tgt.Header = map[string][]string{
				"Content-Type": {"application/json"},
			}
			return nil
		}
	}

	var metrics vegeta.Metrics

	for res := range attacker.Attack(targeter(), rate, duration, "Auth Test") {
		metrics.Add(res)
	}

	metrics.Close()

	// Вывод метрик.
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

	// Проверка SLI времени ответа.
	if metrics.Latencies.Mean > 50*time.Millisecond {
		t.Error("❌ SLI времени ответа не выполнен (больше 50 мс)")
	} else {
		t.Log("✅ SLI времени ответа выполнен")
	}
}
