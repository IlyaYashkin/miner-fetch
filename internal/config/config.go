package config

import (
	"flag"
	"strings"
)

type Config struct {
	TgAPIKey         string   `json:"tg_api_key"`
	TgAdminUsernames []string `json:"tg_admin_usernames"`
	Mode             string   `json:"mode"`
	IsScanner        bool     `json:"is_scanner"`
	Port             string   `json:"port"`
	ParentAuthority  string   `json:"parent_authority"`
}

func GetConfig() Config {
	tgApiKey := flag.String("tg-api-key", "", "Telegram API Key")
	tgAdminUsernames := flag.String("tg-admin-usernames", "", "Telegram admin usernames")
	mode := flag.String("mode", "parent", "Mode")
	noScan := flag.Bool("no-scan", false, "Scan mode")
	port := flag.String("port", "8080", "Port")
	parentAuthority := flag.String("parent-authority", "", "Parent authority")

	flag.Parse()

	return Config{
		TgAPIKey:         *tgApiKey,
		TgAdminUsernames: prepareUsernames(*tgAdminUsernames),
		Mode:             *mode,
		IsScanner:        !*noScan,
		Port:             *port,
		ParentAuthority:  *parentAuthority,
	}
}

func prepareUsernames(usernamesString string) []string {
	usernames := strings.Split(usernamesString, ",")

	for i, username := range usernames {
		usernames[i] = strings.TrimLeft(username, "@")
	}

	return usernames
}
