package main

import (
	"fmt"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"math/rand/v2"
	"time"
)

func generateUsername() string {
	return fmt.Sprintf("user_%d", rand.IntN(10000000))
}

func main() {
	rate := vegeta.Rate{Freq: 200, Per: time.Second}
	duration := 10 * time.Second

	attacker := vegeta.NewAttacker()

	targeter := func() vegeta.Targeter {
		return vegeta.NewStaticTargeter(vegeta.Target{
			Method: "POST",
			URL:    "http://localhost:8080/api/auth",
			Body:   []byte(fmt.Sprintf(`{"username": "%s", "password": "test_pass"}`, generateUsername())),
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
		})
	}

	var metrics vegeta.Metrics

	// Одна атака на всю длительность теста
	for res := range attacker.Attack(targeter(), rate, duration, "Auth Test") {
		metrics.Add(res)
	}

	metrics.Close()

	fmt.Printf("Requests: %d\n", metrics.Requests)
	fmt.Printf("Succes: %.2f%%\n", metrics.Success*100)
	fmt.Printf("Latency (mean): %s\n", metrics.Latencies.Mean)
	fmt.Printf("Status codes: %v\n", metrics.StatusCodes)

	if metrics.Success < 0.9999 {
		fmt.Println("❌ SLI успешности не выполнен (меньше 99.99%)")
	} else {
		fmt.Println("✅ SLI успешности выполнен")
	}

	if metrics.Latencies.Mean > 50*time.Millisecond {
		fmt.Println("❌ SLI времени ответа не выполнен (больше 50 мс)")
	} else {
		fmt.Println("✅ SLI времени ответа выполнен")
	}
}

//Requests: 10000
//Succes: 3.84%
//Latency (mean): 26.057686883s
//Status codes: map[0:9580 200:384 500:36]

// расширен пул соединений 20 => 100
// поставлен таймаут на соединение к бд 100ms

//Requests: 10000
//Succes: 10.04%
//Latency (mean): 828.971533ms
//Status codes: map[200:1004 500:8996]

// расширен пул соединений 100 => 200

//Requests: 10000
//Succes: 10.62%
//Latency (mean): 680.163014ms
//Status codes: map[200:1062 500:8938]

// пул соединений уменьшен до 50

//Requests: 10000
//Succes: 9.09%
//Latency (mean): 474.168418ms
//Status codes: map[200:909 500:9091]

// пул соединений возвращен к 100

//Requests: 10000
//Succes: 11.08%
//Latency (mean): 843.859814ms
//Status codes: map[200:1108 500:8892]
