package scripts

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadTokensFromFile(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Ошибка чтения файла с токенами: %w", err)
	}

	var tokens []string
	if err = json.Unmarshal(data, &tokens); err != nil {
		return nil, fmt.Errorf("Ошибка десериализации токенов: %w", err)
	}

	return tokens, nil
}
