//go:build buy

package load

import (
	"fmt"
	"github.com/magmaheat/merchStore/tests/load_test/scripts"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"golang.org/x/exp/rand"
	"testing"
	"time"
)

func TestBuyLoad(t *testing.T) {
	tokens, err := scripts.LoadTokensFromFile("./register_users/tokens.json")
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
				"Content-Type":  {"application/json"},
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
