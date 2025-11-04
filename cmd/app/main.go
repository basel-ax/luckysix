package main

import (
	"log"

	"github.com/basel-ax/luckysix/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	a, err := app.New()
	err = a.Run()
	if err != nil {
		log.Fatal(err)
	}
}
