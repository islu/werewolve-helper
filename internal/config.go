package internal

type BotConfig struct {
	LineChannelSecret string
	LineChannelToken  string
	Port              string
	DiscordBotToken   string
	DiscordChannelID  string
	LiffID            string

	// DeveloperID     string // Deprecated: developer ID is not used
	// LineNotifyToken string // Deprecated: LINE Notify token is not used
}
