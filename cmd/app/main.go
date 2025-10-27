package main

import (
	"log"

	"github.com/basel-ax/luckysix/config"
	"github.com/basel-ax/luckysix/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
