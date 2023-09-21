package app

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go-tgbot/internal/config"
	"log"
	"regexp"
	"strconv"
	"time"
)

type Task struct {
	Command  string
	Username string
	Date     string
	Text     string
}

func Run() {
	cfg, err := config.Init()
	if err != nil {
		panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		log.Printf(strconv.FormatInt(update.Message.Chat.ID, 10))
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "help":
			msg.Text = "I understand /task"
		case "task":
			log.Printf("[%s] %s", update.Message, update.Message.Text)

			pattern := regexp.MustCompile(`/task\s+@(\S+)\s+(\S+)\s+(.+)$`)
			matches := pattern.FindStringSubmatch(update.Message.Text)

			if len(matches) == 4 {
				command := "task" // Assuming "task" is the default command
				username := matches[1]
				date := matches[2]
				text := matches[3]

				parsedDate, err := time.Parse("02.01.2006", date)
				if err != nil {
					currentYear := time.Now().Year()
					parsedDate, _ = time.Parse("02.01.2006", fmt.Sprintf("%s.%d", date, currentYear))
				}
				formattedDate := parsedDate.Format("02.01.2006")
				task := Task{
					Command:  command,
					Username: username,
					Date:     formattedDate,
					Text:     text,
				}

				msg.Text = task.Text + " " + task.Date
			} else {
				msg.Text = "Invalid input format"
			}

		case "status":
			msg.Text = "I'm ok."
		default:
			msg.Text = "I don't know that command"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
