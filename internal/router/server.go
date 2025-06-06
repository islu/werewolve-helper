package router

import (
	"log"
	"net/http"
	"os"
	"time"
	"werewolve-helper/internal"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"

	_ "github.com/joho/godotenv/autoload"
)

// Initialize the server and start listening for incoming requests.
// This function will set up the necessary routes and handlers.
func StartServer() {
	config := initBotConfig()

	bot, err := messaging_api.NewMessagingApiAPI(config.LineChannelToken)
	if err != nil {
		log.Fatalln(err)
	}

	// Register webhook
	RegisterWebhook(config, bot)
	// Register LIFF page
	RegisterLIFF(config)
	// Register health check
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := &http.Server{
		Addr:         ":" + config.Port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Server starting on port " + config.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %v", err)
	}
}

func initBotConfig() internal.BotConfig {
	channelSecret := mustGetenv("LINE_CHANNEL_SECRET")
	channelToken := mustGetenv("LINE_CHANNEL_TOKEN")
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
		LiffID:            liffID,
		DiscordBotToken:   dcBotToken,
		DiscordChannelID:  dcChannelID,
	}
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Fatal Error: %s environment variable not set.\n", k)
	}
	return v
}
