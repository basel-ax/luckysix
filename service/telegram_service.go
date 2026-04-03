package service

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

// TelegramNotification sends a message to Telegram
func TelegramNotification(message string) error {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	if botToken == "" || chatID == "" {
		log.Println("Telegram credentials not configured. Skipping notification.")
		return nil
	}

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	data := url.Values{}
	data.Set("chat_id", chatID)
	data.Set("text", message)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.PostForm(apiURL, data)
	if err != nil {
		log.Printf("Failed to send Telegram notification: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Telegram API returned status: %d", resp.StatusCode)
		return fmt.Errorf("Telegram API returned status: %d", resp.StatusCode)
	}

	log.Println("Telegram notification sent successfully")
	return nil
}

// SendGenerationNotification sends a notification about generation completion
func SendGenerationNotification(command string, rowsCreated int64, err error) error {
	var message string

	if err != nil {
		message = fmt.Sprintf("❌ %s generation completed with error: %v", command, err)
	} else if rowsCreated == -1 {
		message = fmt.Sprintf("⏭️ %s skipped — already running", command)
	} else {
		message = fmt.Sprintf("✅ %s generation completed. Rows created/updated: %d", command, rowsCreated)
	}

	return TelegramNotification(message)
}
