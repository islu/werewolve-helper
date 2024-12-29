package internal

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/islu/werewolve-helper/internal/notify"
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
	EventCreate = "create"
	EventLook   = "look"
	EventAgain  = "again"
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
				case webhook.ImageMessageContent:
					switch source := e.Source.(type) {
					case webhook.UserSource:
						if err := handleImage(bot, e.ReplyToken, &message, source); err != nil {
							log.Println("Handle image event error: ", err)
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
					if err := handlePostbackEvent(bot, e.ReplyToken, e.Postback, source, config.LIFFID); err != nil {
						log.Println("Handle postback event error: ", err)
					}
				default:
					log.Printf("Unsupported source content: %T\n", e.Source)
				}
			case webhook.FollowEvent:
				switch source := e.Source.(type) {
				case webhook.UserSource:
					if err := push(bot, "FollowEvent", source, config.DiscordBotToken, config.DiscordChannelID); err != nil {
						log.Println("Notify error: ", err)
					}
				default:
					log.Printf("Unsupported source content: %T\n", e.Source)
				}
			case webhook.UnfollowEvent:
				switch source := e.Source.(type) {
				case webhook.UserSource:
					if err := push(bot, "UnfollowEvent", source, config.DiscordBotToken, config.DiscordChannelID); err != nil {
						log.Println("Notify error: ", err)
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

	if ownerID := findRoundByInviteNo(text); ownerID != "" {

		if r, ok := rounds[source.UserId]; ok {

			if r.IsExpired() {
				delete(rounds, ownerID)
				m1 := messaging_api.TextMessage{Text: "活動已結束"}
				return reply(bot, replyToken, m1)
			}

			user, err := bot.GetProfile(source.UserId)
			if err != nil {
				return err
			}

			if ok, p := r.IsRegistrationDuplicate(user.UserId); ok {
				m1 := messaging_api.TextMessage{Text: "已註冊，你的身分是 " + p.Identity.String()}
				return reply(bot, replyToken, m1)
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

		m1 := messaging_api.TextMessage{Text: "查無此活動"}
		return reply(bot, replyToken, m1)
	}

	return errors.New("Unknown message text " + text)
}

func handleImage(bot *messaging_api.MessagingApiAPI, replyToken string, message *webhook.ImageMessageContent, source webhook.UserSource) error {

	u := message.ContentProvider.OriginalContentUrl

	url, err := url.Parse(u)
	if err != nil {
		return err
	}
	q := url.Query()

	switch q.Get("m") {
	case "settingRole":

		// Generate inviteNo
		inviteNo := fmt.Sprintf("%06d", rand.Intn(999999))
		if isRoundInviteNoDuplicate(inviteNo) {
			log.Println("inviteNo duplicate: " + inviteNo)
			m1 := messaging_api.TextMessage{Text: "創建失敗，請重新嘗試"}
			return reply(bot, replyToken, m1)
		}
		// Create round and set identity
		round := NewRound(source.UserId, inviteNo)
		if v := q.Get("b1"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				log.Println("parse error with b1: ", err)
				return err
			}
			round.SetIdentity(source.UserId, WerewolfKing, n)
		}
		if v := q.Get("b2"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				log.Println("parse error with b2: ", err)
				return err
			}
			round.SetIdentity(source.UserId, WhiteWerewolf, n)
		}
		if v := q.Get("b3"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				log.Println("parse error with b3: ", err)
				return err
			}
			round.SetIdentity(source.UserId, GhostRider, n)
		}
		if v := q.Get("b4"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				log.Println("parse error with b4: ", err)
				return err
			}
			round.SetIdentity(source.UserId, WerewolfBeauty, n)
		}
		if v := q.Get("b0"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				log.Println("parse error with b0: ", err)
				return err
			}
			round.SetIdentity(source.UserId, Werewolf, n)
		}
		if v := q.Get("g1"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				log.Println("parse error with g1: ", err)
				return err
			}
			round.SetIdentity(source.UserId, Seer, n)
		}
		if v := q.Get("g2"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				log.Println("parse error with g2: ", err)
				return err
			}
			round.SetIdentity(source.UserId, Witch, n)
		}
		if v := q.Get("g3"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				log.Println("parse error with g3: ", err)
				return err
			}
			round.SetIdentity(source.UserId, Hunter, n)
		}
		if v := q.Get("g4"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				log.Println("parse error with g4: ", err)
				return err
			}
			round.SetIdentity(source.UserId, Guard, n)
		}
		if v := q.Get("g5"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				log.Println("parse error with g5: ", err)
				return err
			}
			round.SetIdentity(source.UserId, Knight, n)
		}
		if v := q.Get("g6"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				log.Println("parse error with g6: ", err)
				return err
			}
			round.SetIdentity(source.UserId, Magician, n)
		}
		if v := q.Get("g0"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				log.Println("parse error with g0: ", err)
				return err
			}
			round.SetIdentity(source.UserId, Villager, n)
		}

		rounds[source.UserId] = round

		m1 := messaging_api.TextMessage{Text: "成功創建房間編號為: " + inviteNo}
		return reply(bot, replyToken, m1)

	}
	return errors.New("Unknown url query key " + q.Get("m"))
}

func handlePostbackEvent(bot *messaging_api.MessagingApiAPI,
	replyToken string,
	postback *webhook.PostbackContent,
	source webhook.UserSource,
	liffID string,
) error {
	switch postback.Data {
	case EventCreate:

		delete(rounds, source.UserId)

		return reply(bot, replyToken, ModeSettingTemplateV2(liffID))

	case EventLook:

		if r, ok := rounds[source.UserId]; ok {
			m1 := messaging_api.TextMessage{Text: "房間編號為: " + r.InviteNo}
			m2 := messaging_api.TextMessage{Text: r.GetParticipantsInfoReplyMessage(source.UserId)}
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
	messages = append(messages, msg...)

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

func push(bot *messaging_api.MessagingApiAPI, eventType string, source webhook.UserSource, discordBotToken, discordChannelID string) error {

	profile, err := bot.GetProfile(source.UserId)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("%s\n\n%#v", eventType, profile)
	_, err = notify.SendMessageByDiscord(discordBotToken, discordChannelID, msg)
	if err != nil {
		return err
	}
	return nil
}
