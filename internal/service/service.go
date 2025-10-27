// Package service implements application business logic. Each logic group in own file.
package service

// Services contains all services
type Services struct {
	Translation Translation
	Tasks       Tasks
	Bip39       Bip39
	Blockchain  Blockchain
	Telegram    Telegram
}
