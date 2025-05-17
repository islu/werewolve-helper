package internal

type BotConfig struct {
	LineChannelSecret string
	LineChannelToken  string
	Port              string
	DeveloperID       string // Duplicate: duplicate field should be removed
	DiscordBotToken   string
	DiscordChannelID  string
	LiffID            string
	LineNotifyToken   string // Duplicate: duplicate field should be removed
}
