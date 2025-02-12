package main

import "github.com/magmaheat/merchStore/internal/app"

const configPath = "config/config.yaml"

func main() {
	app.Run(configPath)
}
