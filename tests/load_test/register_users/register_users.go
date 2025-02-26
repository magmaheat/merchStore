package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	apiURL      = "http://localhost:8080/api/auth"
	outputFile  = "./tests/load_test/register_users/tokens.json"
	userCount   = 100000
	workerLimit = 100
)

func main() {
	tokens, err := registerUsers(userCount)
	if err != nil {
		log.Fatalf("Ошибка при регистрации пользователей: %v", err)
	}

	err = saveTokensToFile(tokens, outputFile)
	if err != nil {
		log.Fatalf("Ошибка при сохранении токенов: %v", err)
	}

	log.Printf("✅ Успешно зарегистрировано %d пользователей. Токены сохранены в %s", userCount, outputFile)
}

func generateUserList(count int) []string {
	users := make([]string, count)
	for i := 0; i < count; i++ {
		users[i] = fmt.Sprintf("user_%d", i+1)
	}
	return users
}

func registerUsers(count int) ([]string, error) {
	startTime := time.Now()
	users := generateUserList(count)
	tokens := make([]string, 0, count)
	var mu sync.Mutex

	sem := make(chan struct{}, workerLimit)
	errCh := make(chan error, count)
	var wg sync.WaitGroup

	log.Info("Начало регистрации пользователей.")

	for idx, user := range users {
		sem <- struct{}{}
		wg.Add(1)

		go func(idx int, user string) {
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

			if idx != 0 && idx%500 == 0 {
				log.Infof("Пользователей зарегистрировано: %d\n", idx)
			}
		}(idx, user)
	}

	wg.Wait()
	close(errCh)

	if len(errCh) > 0 {
		return nil, <-errCh
	}

	log.Infof("⏳ Время регистрации пользователей: %v", time.Since(startTime))
	return tokens, nil
}

func registerUser(user string) (string, error) {
	resp, err := http.Post(
		apiURL,
		"application/json",
		strings.NewReader(fmt.Sprintf(`{"username": "%s", "password": "test_pass"}`, user)),
	)
	if err != nil {
		return "", fmt.Errorf("ошибка регистрации %s: %w", user, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ошибка регистрации %s: статус %d", user, resp.StatusCode)
	}

	var result struct {
		Token string `json:"token"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("ошибка декодирования ответа для %s: %w", user, err)
	}

	return result.Token, nil
}

func saveTokensToFile(tokens []string, filename string) error {
	data, err := json.Marshal(tokens)
	if err != nil {
		return fmt.Errorf("ошибка сериализации токенов: %w", err)
	}

	return os.WriteFile(filename, data, 0644)
}
