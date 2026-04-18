package handler

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler func(*tgbotapi.Message) error

type BotHandler struct {
	bot      *tgbotapi.BotAPI
	commands map[string]CommandHandler
	service  Service
}

func NewBotHandler(bot *tgbotapi.BotAPI, service Service) *BotHandler {
	h := &BotHandler{
		bot:      bot,
		commands: make(map[string]CommandHandler),
		service:  service,
	}

	h.RegisterCommand("start", h.handleStart)
	h.RegisterCommand("help", h.handleHelp)
	h.RegisterCommand("boost", h.handleBoost)

	return h
}

func (h *BotHandler) RegisterCommand(command string, handler CommandHandler) {
	h.commands[strings.ToLower(command)] = handler
}

func (h *BotHandler) Run(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if err := h.handleMessage(update.Message); err != nil {
			log.Printf("handle message error: %v", err)
		}
	}
}

func (h *BotHandler) handleMessage(msg *tgbotapi.Message) error {
	if !msg.IsCommand() {
		return h.reply(msg.Chat.ID, "Используй /help чтобы увидеть доступные команды.", msg.MessageID)
	}

	command := strings.ToLower(msg.Command())
	handler, ok := h.commands[command]
	if !ok {
		return h.reply(msg.Chat.ID, "Неизвестная команда. Используй /help.", msg.MessageID)
	}

	return handler(msg)
}

func (h *BotHandler) handleStart(msg *tgbotapi.Message) error {
	text := "Привет! Я групповой бот прокачки статы.\n\n" +
		"Используй /boost раз в 24 часа, чтобы попытаться увеличить показатель.\n" +
		"Если использовать раньше времени, можно получить штраф.\n\n" +
		"Команды:\n" +
		"/boost - попытка прокачки\n" +
		"/help - показать команды"

	return h.reply(msg.Chat.ID, text, msg.MessageID)
}

func (h *BotHandler) handleBoost(msg *tgbotapi.Message) error {
	result, err := h.service.Boost(msg.From.ID, msg.Chat.ID)
	if err != nil {
		log.Printf("boost error: %v", err)
		return h.reply(msg.Chat.ID, "Произошла ошибка, попробуй позже.", msg.MessageID)
	}

	name := displayName(msg)
	changeLine := fmt.Sprintf("%s, твой писюн вырос на *%d* см.", name, result.Delta)
	if result.Delta < 0 {
		changeLine = fmt.Sprintf("%s, твой писюн уменьшился на *%d* см.", name, abs64(result.Delta))
	}

	if !result.Penalized {

		nextAttemptLine := "Следующая попытка завтра!"

		text := fmt.Sprintf(
			"%s\nТеперь он равен *%d* см.\nТы занимаешь *%d* %s в топе.\n%s",
			changeLine,
			result.Length,
			result.Rank,
			placeWord(result.Rank),
			nextAttemptLine,
		)

		return h.reply(msg.Chat.ID, text, msg.MessageID)
	} else {
		penalizedMessage := fmt.Sprintf("%s, к сожалению ты поспешил, и будешь наказан(", name)
		changeLine = fmt.Sprintf("Твой писюн уменьшился на *%d* см.", abs64(result.Delta))
		text := fmt.Sprintf(
			"%s\n%s\nТеперь он равен *%d* см.\nТы занимаешь *%d* %s в топе.",
			penalizedMessage,
			changeLine,
			result.Length,
			result.Rank,
			placeWord(result.Rank),
		)
		return h.reply(msg.Chat.ID, text, msg.MessageID)
	}

}

func (h *BotHandler) handleHelp(msg *tgbotapi.Message) error {
	text := strings.Join([]string{
		"Доступные команды:",
		"/start - приветствие",
		"/boost - попытка прокачки",
		"/help - показать команды",
	}, "\n")
	return h.reply(msg.Chat.ID, text, msg.MessageID)
}

func (h *BotHandler) reply(chatID int64, text string, replyTo int) error {
	message := tgbotapi.NewMessage(chatID, text)
	message.ParseMode = "MARKDOWN"
	message.ReplyToMessageID = replyTo
	_, err := h.bot.Send(message)
	return err
}

func displayName(msg *tgbotapi.Message) string {
	if msg.From.FirstName != "" {
		return msg.From.FirstName
	}
	if msg.From.UserName != "" {
		return msg.From.UserName
	}
	return "Игрок"
}

func abs64(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}

func placeWord(rank int) string {
	mod10 := rank % 10
	mod100 := rank % 100

	if mod10 == 1 && mod100 != 11 {
		return "место"
	}
	if mod10 >= 2 && mod10 <= 4 && (mod100 < 12 || mod100 > 14) {
		return "места"
	}
	return "мест"
}
