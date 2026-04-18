package main

import (
	"PipisaBot/internal/handler"
	"PipisaBot/internal/repository"
	"PipisaBot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Println("TOKEN is empty")
		os.Exit(1)
	}

	webhookBaseURL := os.Getenv("WEBHOOK_BASE_URL")
	if webhookBaseURL == "" {
		log.Println("WEBHOOK_BASE_URL is empty")
		os.Exit(1)
	}

	listenAddr := os.Getenv("WEBHOOK_LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":8080"
	}

	webhookPath := os.Getenv("WEBHOOK_PATH")
	if webhookPath == "" {
		webhookPath = "/webhook"
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	webhookPath = normalizeWebhookPath(webhookPath) + "/" + token
	webhookURL := strings.TrimRight(webhookBaseURL, "/") + webhookPath

	updates := bot.ListenForWebhook(webhookPath)

	go func() {
		log.Printf("webhook server listen on %s%s", listenAddr, webhookPath)
		if err := http.ListenAndServe(listenAddr, nil); err != nil {
			log.Printf("webhook server error: %v", err)
			os.Exit(1)
		}
	}()

	webhookConfig, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	webhookConfig.DropPendingUpdates = true

	if _, err = bot.Request(webhookConfig); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Printf("webhook configured: %s", webhookURL)

	repo := repository.NewUserRepository()
	userService := service.NewUserService(repo)

	botHandler := handler.NewBotHandler(bot, userService)
	botHandler.Run(updates)
}

func normalizeWebhookPath(path string) string {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return "/webhook"
	}
	if !strings.HasPrefix(trimmed, "/") {
		return "/" + trimmed
	}
	return trimmed
}
