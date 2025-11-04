package service

import (
	"context"
	"log"
)

// LogNotifier -.
type LogNotifier struct{}

// NewLogNotifier -.
func NewLogNotifier() *LogNotifier {
	return &LogNotifier{}
}

// Notify -.
func (n *LogNotifier) Notify(ctx context.Context, message string) error {
	log.Println(message)
	return nil
}
