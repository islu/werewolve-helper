package notify

import (
	"github.com/bwmarrin/discordgo"
)

func SendMessageByDiscord(accessToken, channelID, message string) (response *discordgo.Message, err error) {

	discord, err := discordgo.New("Bot " + accessToken)
	if err != nil {
		return nil, err
	}

	msg, err := discord.ChannelMessageSend(channelID, message)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
