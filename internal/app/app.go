package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go-tgbot/internal/config"
	"log"
	"net/http"
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
				err = sendService(&task, cfg.Service.Url, cfg.Service.Username, cfg.Service.Password)
				if err == nil {
					msg.Text = "Task sent successfully"
				} else {
					msg.Text = "Failed to send task"
				}
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

func sendService(task *Task, url string, username, password string) error {
	client := &http.Client{}
	jsonData, err := json.Marshal(task)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Task data sent successfully.")
		return nil
	} else {
		fmt.Println("Failed to send task data. Status code:", resp.StatusCode)
		return fmt.Errorf("Failed to send task data. Status code: %d", resp.StatusCode)
	}
}
