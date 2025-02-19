package config

import (
	"flag"
	"strings"
)

type Config struct {
	TelegramAPIKey         string
	TelegramAdminUsernames []string
}

func GetConfig() Config {
	telegramApiKey := flag.String("telegram-api-key", "", "Telegram API Key")
	telegramAdminUsernames := flag.String("telegram-admin-usernames", "", "Telegram admin usernames")
	flag.Parse()

	return Config{
		TelegramAPIKey:         *telegramApiKey,
		TelegramAdminUsernames: prepareUsernames(*telegramAdminUsernames),
	}
}

func prepareUsernames(usernamesString string) []string {
	usernames := strings.Split(usernamesString, ",")

	for i, username := range usernames {
		usernames[i] = strings.TrimLeft(username, "@")
	}

	return usernames
}
