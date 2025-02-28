package config

import (
	"os"
	"strings"
)

type Config struct {
	TgAPIKey         string   `json:"tg_api_key"`
	TgAdminUsernames []string `json:"tg_admin_usernames"`
	Mode             string   `json:"mode"`
	IsScanner        bool     `json:"is_scanner"`
	Port             string   `json:"port"`
	ParentAuthority  string   `json:"parent_authority"`
	NodeName         string   `json:"node_name"`
	AuthKey          string   `json:"auth_key"`
}

func GetConfig() Config {
	tgApiKey := getEnv("TG_API_KEY", "")
	tgAdminUsernames := getEnv("TG_ADMIN_USERNAMES", "")
	mode := getEnv("MODE", "parent")
	isScanner := getEnv("IS_SCANNER", "true")
	port := getEnv("PORT", "8080")
	parentAuthority := getEnv("PARENT_AUTHORITY", "")
	nodeName := getEnv("NODE_NAME", "Unknown")
	authKey := getEnv("AUTH_KEY", "")

	return Config{
		TgAPIKey:         tgApiKey,
		TgAdminUsernames: prepareUsernames(tgAdminUsernames),
		Mode:             mode,
		IsScanner:        isScanner == "true",
		Port:             port,
		ParentAuthority:  parentAuthority,
		NodeName:         nodeName,
		AuthKey:          authKey,
	}
}

func prepareUsernames(usernamesString string) []string {
	usernames := strings.Split(usernamesString, ",")

	for i, username := range usernames {
		usernames[i] = strings.TrimLeft(username, "@")
	}

	return usernames
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
