package main

import (
	"log"
	"net/http"
	"os"

	"github.com/islu/werewolve-helper/internal"

	_ "github.com/joho/godotenv/autoload"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

func main() {

	config := initBotConfig()

	bot, err := messaging_api.NewMessagingApiAPI(config.LineChannelToken)
	if err != nil {
		log.Fatalln(err)
	}

	internal.RegisterWebhook(config, bot)
	// Add LIFF page endpoint
	internal.RegisterLIFF(config)
	// Add health check endpoint
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatalln(err)
	}

	log.Println("Server start running")
}

func initBotConfig() internal.BotConfig {
	channelSecret := mustGetenv("LINE_CHANNEL_SECRET")
	channelToken := mustGetenv("LINE_CHANNEL_TOKEN")
	developerID := mustGetenv("DEVELOPER_ID")
	liffID := mustGetenv("LIFF_ID")
	dcBotToken := mustGetenv("DISCORD_BOT_TOKEN")
	dcChannelID := mustGetenv("DISCORD_CHANNEL_ID")

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	return internal.BotConfig{
		LineChannelSecret: channelSecret,
		LineChannelToken:  channelToken,
		Port:              port,
		DeveloperID:       developerID,
		LIFFID:            liffID,
		DiscordBotToken:   dcBotToken,
		DiscordChannelID:  dcChannelID,

		// LineNotifyToken:   notifyToken,
	}
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Fatal Error in connect_connector.go: %s environment variable not set.\n", k)
	}
	return v
}
