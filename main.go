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
		log.Fatal(err)
	}

	internal.RegisterWebhook(config, bot)

	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatal(err)
	}
}

func initBotConfig() internal.BotConfig {
	channelSecret := mustGetenv("LINE_CHANNEL_SECRET")
	channelToken := mustGetenv("LINE_CHANNEL_TOKEN")
	developerID := mustGetenv("DEVELOPER_ID")

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	return internal.BotConfig{
		LineChannelSecret: channelSecret,
		LineChannelToken:  channelToken,
		Port:              port,
		DeveloperID:       developerID,
	}
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Fatal Error in connect_connector.go: %s environment variable not set.\n", k)
	}
	return v
}
