package config

import "os"

var (
	ServerURL = os.Getenv("SERVER_URL")
	BotSocket = os.Getenv("BOT_SOCKET")
	DbSocket  = os.Getenv("DB_SOCKET")
	Mode      = os.Getenv("MODE")
	BotToken  = os.Getenv("BOT_TOKEN")
	TeamName  = os.Getenv("TEAM_NAME")
)
