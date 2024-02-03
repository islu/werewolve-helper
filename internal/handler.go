package internal

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/line/line-bot-sdk-go/v8/linebot"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

var (
	// {key: ownerID, value: Round}
	rounds = make(map[string]*Round)
)

// Postback event key
const (
	EventCreate              = "create"
	EventLook                = "look"
	EventAgain               = "again"
	Event9PersonStandardMode = "9人標準配置"
	EventCustomMode          = "自訂配置"
	EventSettingWerewolf     = "設定狼人"
	EventSettingVillager     = "設定平民"
	EventSettingSeer         = "設定預言家"
	EventSettingWitch        = "設定女巫"
	EventSettingHunter       = "設定獵人"
)

func RegisterWebhook(config BotConfig, bot *messaging_api.MessagingApiAPI) {

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		// log.Println("/callback called...")

		cb, err := webhook.ParseRequest(config.LineChannelSecret, req)
		if err != nil {
			log.Printf("Cannot parse request: %+v\n", err)
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}

		// log.Println("Handling events...")
		for _, event := range cb.Events {
			// log.Printf("/callback called%+v...\n", event)

			switch e := event.(type) {
			case webhook.MessageEvent:
				switch message := e.Message.(type) {
				case webhook.TextMessageContent:
					switch source := e.Source.(type) {
					case webhook.UserSource:
						if err := handleText(bot, e.ReplyToken, &message, source); err != nil {
							log.Println("Handle text event error: ", err)
						}
					default:
						log.Printf("Unsupported source content: %T\n", e.Source)
					}
				default:
					log.Printf("Unsupported message content: %T\n", e.Message)
				}
			case webhook.PostbackEvent:
				switch source := e.Source.(type) {
				case webhook.UserSource:
					if err := handlePostbackEvent(bot, e.ReplyToken, e.Postback, source); err != nil {
						log.Println("Handle postback event error: ", err)
					}
				default:
					log.Printf("Unsupported source content: %T\n", e.Source)
				}
			default:
				log.Printf("Unsupported event: %T\n", event)
			}
		}
	})
}

func handleText(bot *messaging_api.MessagingApiAPI, replyToken string, message *webhook.TextMessageContent, source webhook.UserSource) error {

	text := message.Text

	switch text {
	case "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13":

		if r, ok := rounds[source.UserId]; ok {
			if r.TempIdentityFlag {
				n, _ := strconv.Atoi(text)
				r.SetIdentity(source.UserId, r.TempIdentity, n)
			}
			r.TempIdentityFlag = false

			m1 := messaging_api.TextMessage{Text: "設定成功"}
			return reply(bot, replyToken, m1)
		}

	}

	if ownerID := findRoundByInviteNo(text); ownerID != "" {

		if r, ok := rounds[source.UserId]; ok {

			user, err := bot.GetProfile(source.UserId)
			if err != nil {
				return err
			}

			iden := r.Register(source.UserId, user.DisplayName, user.PictureUrl)
			if iden == "" {
				m1 := messaging_api.TextMessage{Text: "已額滿"}
				return reply(bot, replyToken, m1)
			}
			var sb strings.Builder
			sb.WriteString("你的身分是 ")
			sb.WriteString(iden)
			m1 := messaging_api.TextMessage{Text: sb.String()}
			return reply(bot, replyToken, m1)
		}
	}

	return errors.New("Unknown message text " + text)
}

func handlePostbackEvent(bot *messaging_api.MessagingApiAPI, replyToken string, postback *webhook.PostbackContent, source webhook.UserSource) error {
	switch postback.Data {
	case EventCreate:

		delete(rounds, source.UserId)
		return reply(bot, replyToken, ModeSettingTemplate())

	case EventLook:

		if r, ok := rounds[source.UserId]; ok {
			m1 := messaging_api.TextMessage{Text: "房間編號為: " + r.InviteNo}
			m2 := messaging_api.TextMessage{Text: r.GetParticipantsInfoStr(source.UserId)}
			return reply(bot, replyToken, m1, m2)
		}

		m1 := messaging_api.TextMessage{Text: "...目前沒有開設房間\n請先開設房間喔"}
		return reply(bot, replyToken, m1)

	case EventAgain:

		if r, ok := rounds[source.UserId]; ok {
			r.Again()
			m1 := messaging_api.TextMessage{Text: "已經重新發牌囉!"}
			return reply(bot, replyToken, m1)
		}

		m1 := messaging_api.TextMessage{Text: "...目前沒有開設房間\n請先開設房間喔"}
		return reply(bot, replyToken, m1)

	case Event9PersonStandardMode:

		inviteNo := fmt.Sprintf("%06d", rand.Intn(999999))
		if isRoundInviteNoDuplicate(inviteNo) {
			log.Println("inviteNo duplicate: " + inviteNo)
			m1 := messaging_api.TextMessage{Text: "創建失敗，請重新嘗試"}
			return reply(bot, replyToken, m1)
		}

		rounds[source.UserId] = NewRoundWith9PersonStandardMode(source.UserId, inviteNo)

		m1 := messaging_api.TextMessage{Text: "成功創建房間編號為: " + inviteNo}
		return reply(bot, replyToken, m1)

	case EventCustomMode:

		inviteNo := fmt.Sprintf("%06d", rand.Intn(999999))
		if isRoundInviteNoDuplicate(inviteNo) {
			log.Println("inviteNo duplicate: " + inviteNo)
			m1 := messaging_api.TextMessage{Text: "創建失敗，請重新嘗試"}
			return reply(bot, replyToken, m1)
		}

		rounds[source.UserId] = NewRound(source.UserId, inviteNo)

		return reply(bot, replyToken, CustomModeTemplate())

	case EventSettingWerewolf:

		if r, ok := rounds[source.UserId]; ok {
			r.TempIdentity = Werewolf
			r.TempIdentityFlag = true
			return reply(bot, replyToken, QuickReplyButtons())
		}
		return nil

	case EventSettingVillager:

		if r, ok := rounds[source.UserId]; ok {
			r.TempIdentity = Villager
			r.TempIdentityFlag = true
			return reply(bot, replyToken, QuickReplyButtons())
		}
		return nil

	case EventSettingSeer:

		if r, ok := rounds[source.UserId]; ok {
			r.TempIdentity = Seer
			r.TempIdentityFlag = true
			return reply(bot, replyToken, QuickReplyButtons())
		}
		return nil

	case EventSettingWitch:

		if r, ok := rounds[source.UserId]; ok {
			r.TempIdentity = Witch
			r.TempIdentityFlag = true
			return reply(bot, replyToken, QuickReplyButtons())
		}
		return nil

	case EventSettingHunter:

		if r, ok := rounds[source.UserId]; ok {
			r.TempIdentity = Hunter
			r.TempIdentityFlag = true
			return reply(bot, replyToken, QuickReplyButtons())
		}
		return nil
	}
	return errors.New("Unknown event key " + postback.Data)
}

func isRoundInviteNoDuplicate(inviteNo string) bool {
	for _, r := range rounds {
		if r.InviteNo == inviteNo {
			return true
		}
	}
	return false
}

// Return round owner ID
func findRoundByInviteNo(inviteNo string) string {
	for id, r := range rounds {
		if r.InviteNo == inviteNo {
			return id
		}
	}
	return ""
}

func reply(bot *messaging_api.MessagingApiAPI, replyToken string, msg ...messaging_api.MessageInterface) error {

	var messages []messaging_api.MessageInterface
	for _, m := range msg {
		messages = append(messages, m)
	}

	if _, err := bot.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: replyToken,
			Messages:   messages,
		},
	); err != nil {
		return err
	}
	return nil
}
