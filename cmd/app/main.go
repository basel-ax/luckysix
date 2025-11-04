package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// a, err := app.New()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = a.Run()
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
