package config

import "github.com/ilyakaznacheev/cleanenv"

type (
	Config struct {
		Telegram TelegramBot `yaml:"telegram_bot"`
		Service  Service1c   `yaml:"service"`
	}
	TelegramBot struct {
		BotToken string `yaml:"bot_token"`
	}
	Service1c struct {
		Url string `yaml:"url"`
	}
)

func Init() (*Config, error) {
	cfg := Config{}
	err := cleanenv.ReadConfig("./config.yml", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
