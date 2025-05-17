package notify

import (
	"github.com/bwmarrin/discordgo"
)

// SendMessageByDiscord sends a message to a Discord channel using a bot token.
// It takes an access token, channel ID, and the message string as input.
// It returns the sent message object or an error if one occurred.
func SendMessageByDiscord(accessToken, channelID, message string) (*discordgo.Message, error) {
	discord, err := discordgo.New("Bot " + accessToken)
	if err != nil {
		return nil, err
	}

	msg, err := discord.ChannelMessageSend(channelID, message)
	return msg, err
}
