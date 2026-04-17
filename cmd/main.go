package main

import (
	"PipisaBot/internal/handler"
	"PipisaBot/internal/repository"
	"PipisaBot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
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

	client, err := newHTTPClientFromEnv()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	bot, err := tgbotapi.NewBotAPIWithClient(token, tgbotapi.APIEndpoint, client)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	repo := repository.NewUserRepository()
	userService := service.NewUserService(repo)

	botHandler := handler.NewBotHandler(bot, userService)
	botHandler.Run()
}

//func newHTTPClientFromEnv() (*http.Client, error) {
//	transport := &http.Transport{
//		Proxy: http.ProxyFromEnvironment,
//	}
//
//	if proxyRaw := os.Getenv("PROXY_URL"); proxyRaw != "" {
//		proxyURL, err := url.Parse(proxyRaw)
//		if err != nil {
//			return nil, err
//		}
//		transport.Proxy = http.ProxyURL(proxyURL)
//		log.Printf("using proxy: %s", proxyURL.Redacted())
//	}
//
//	return &http.Client{
//		Timeout:   20 * time.Second,
//		Transport: transport,
//	}, nil
//}
