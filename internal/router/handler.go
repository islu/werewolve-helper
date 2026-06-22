package router

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"werewolve-helper/internal"
	"werewolve-helper/internal/adapter/notify"
	"werewolve-helper/internal/domain"
	"werewolve-helper/internal/usecase"

	"github.com/line/line-bot-sdk-go/v8/linebot"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

// Postback event key
const (
	EventCreate = "create"
	EventLook   = "look"
	EventAgain  = "again"
)

// roleKeys maps the LIFF image query keys to game identities. Used when a round
// is created so new roles only require adding a row here.
var roleKeys = []struct {
	key      string
	identity domain.Identity
}{
	{"b1", domain.WerewolfKing},
	{"b2", domain.WhiteWerewolf},
	{"b3", domain.GhostRider},
	{"b4", domain.WerewolfBeauty},
	{"b0", domain.Werewolf},
	{"g1", domain.Seer},
	{"g2", domain.Witch},
	{"g3", domain.Hunter},
	{"g4", domain.Guard},
	{"g5", domain.Knight},
	{"g6", domain.Magician},
	{"g0", domain.Villager},
}

func RegisterWebhook(config internal.BotConfig, bot *messaging_api.MessagingApiAPI) {
	manager := usecase.NewRoundManager()

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		cb, err := webhook.ParseRequest(config.LineChannelSecret, req)
		if err != nil {
			log.Printf("Cannot parse request: %+v\n", err)
			if errors.Is(err, linebot.ErrInvalidSignature) {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}

		for _, event := range cb.Events {
			switch e := event.(type) {
			case webhook.MessageEvent:
				source, ok := userSource(e.Source)
				if !ok {
					continue
				}
				switch message := e.Message.(type) {
				case webhook.TextMessageContent:
					if err := handleText(bot, manager, e.ReplyToken, &message, source); err != nil {
						log.Println("Handle text event error: ", err)
					}
				case webhook.ImageMessageContent:
					if err := handleImage(bot, manager, e.ReplyToken, &message, source); err != nil {
						log.Println("Handle image event error: ", err)
					}
				default:
					log.Printf("Unsupported message content: %T\n", e.Message)
				}
			case webhook.PostbackEvent:
				source, ok := userSource(e.Source)
				if !ok {
					continue
				}
				if err := handlePostbackEvent(bot, manager, e.ReplyToken, e.Postback, source, config.LiffID); err != nil {
					log.Println("Handle postback event error: ", err)
				}
			case webhook.FollowEvent:
				source, ok := userSource(e.Source)
				if !ok {
					continue
				}
				if err := push(bot, "FollowEvent", source, config.DiscordBotToken, config.DiscordChannelID); err != nil {
					log.Println("Notify error: ", err)
				}
			case webhook.UnfollowEvent:
				source, ok := userSource(e.Source)
				if !ok {
					continue
				}
				if err := push(bot, "UnfollowEvent", source, config.DiscordBotToken, config.DiscordChannelID); err != nil {
					log.Println("Notify error: ", err)
				}
			default:
				log.Printf("Unsupported event: %T\n", event)
			}
		}
	})
}

// userSource extracts a UserSource, logging when the source type is unsupported.
func userSource(src webhook.SourceInterface) (webhook.UserSource, bool) {
	s, ok := src.(webhook.UserSource)
	if !ok {
		log.Printf("Unsupported source content: %T\n", src)
		return webhook.UserSource{}, false
	}
	return s, true
}

func handleText(bot *messaging_api.MessagingApiAPI, manager *usecase.RoundManager, replyToken string, message *webhook.TextMessageContent, source webhook.UserSource) error {
	text := message.Text

	// The text is treated as an invite number; find the round it belongs to.
	r, ok := manager.FindByInviteNo(text)
	if !ok {
		return errors.New("Unknown message text " + text)
	}

	if r.IsExpired() {
		manager.Delete(r.OwnerID)
		return reply(bot, replyToken, messaging_api.TextMessage{Text: "活動已結束"})
	}

	user, err := bot.GetProfile(source.UserId)
	if err != nil {
		return err
	}

	if dup, p := r.IsRegistrationDuplicate(user.UserId); dup {
		return reply(bot, replyToken, messaging_api.TextMessage{Text: "已註冊，你的身分是 " + p.Identity.String()})
	}

	iden := r.Register(source.UserId, user.DisplayName, user.PictureUrl)
	if iden == "" {
		return reply(bot, replyToken, messaging_api.TextMessage{Text: "已額滿"})
	}
	return reply(bot, replyToken, messaging_api.TextMessage{Text: "你的身分是 " + iden})
}

func handleImage(bot *messaging_api.MessagingApiAPI, manager *usecase.RoundManager, replyToken string, message *webhook.ImageMessageContent, source webhook.UserSource) error {
	parsed, err := url.Parse(message.ContentProvider.OriginalContentUrl)
	if err != nil {
		return err
	}
	q := parsed.Query()

	switch q.Get("m") {
	case "settingRole":
		// Generate inviteNo
		randomNo, err := domain.Rng.IntN(999999)
		if err != nil {
			log.Println("Error generating random inviteNo: ", err)
			return err
		}
		inviteNo := fmt.Sprintf("%06d", randomNo)
		if manager.IsInviteNoDuplicate(inviteNo) {
			log.Println("inviteNo duplicate: " + inviteNo)
			return reply(bot, replyToken, messaging_api.TextMessage{Text: "創建失敗，請重新嘗試"})
		}

		// Create round and set identities from the query keys.
		round := domain.NewRound(source.UserId, inviteNo)
		for _, rk := range roleKeys {
			v := q.Get(rk.key)
			if v == "" {
				continue
			}
			n, err := strconv.Atoi(v)
			if err != nil {
				log.Printf("parse error with %s: %v", rk.key, err)
				return err
			}
			round.SetIdentity(source.UserId, rk.identity, n)
		}

		manager.Put(round)
		return reply(bot, replyToken, messaging_api.TextMessage{Text: "成功創建房間編號為: " + inviteNo})
	}
	return errors.New("Unknown url query key " + q.Get("m"))
}

func handlePostbackEvent(bot *messaging_api.MessagingApiAPI,
	manager *usecase.RoundManager,
	replyToken string,
	postback *webhook.PostbackContent,
	source webhook.UserSource,
	liffID string,
) error {
	switch postback.Data {
	case EventCreate:
		manager.Delete(source.UserId)
		return reply(bot, replyToken, ModeSettingTemplateV2(liffID))

	case EventLook:
		if r, ok := manager.Get(source.UserId); ok {
			m1 := messaging_api.TextMessage{Text: "房間編號為: " + r.InviteNo}
			m2 := messaging_api.TextMessage{Text: r.GetParticipantsInfoReplyMessage(source.UserId)}
			return reply(bot, replyToken, m1, m2)
		}
		return reply(bot, replyToken, messaging_api.TextMessage{Text: "...目前沒有開設房間\n請先開設房間喔"})

	case EventAgain:
		if r, ok := manager.Get(source.UserId); ok {
			r.Again()
			return reply(bot, replyToken, messaging_api.TextMessage{Text: "已經重新發牌囉!"})
		}
		return reply(bot, replyToken, messaging_api.TextMessage{Text: "...目前沒有開設房間\n請先開設房間喔"})
	}

	return errors.New("Unknown event key " + postback.Data)
}

func reply(bot *messaging_api.MessagingApiAPI, replyToken string, msg ...messaging_api.MessageInterface) error {
	if _, err := bot.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: replyToken,
			Messages:   msg,
		},
	); err != nil {
		return err
	}
	return nil
}

func push(bot *messaging_api.MessagingApiAPI, eventType string, source webhook.UserSource, discordBotToken, discordChannelID string) error {
	profile, err := bot.GetProfile(source.UserId)
	if err != nil {
		profile = &messaging_api.UserProfileResponse{}
		// return err
	}

	msg := fmt.Sprintf("%s\n\n%#v", eventType, profile)
	_, err = notify.SendMessageByDiscord(discordBotToken, discordChannelID, msg)
	if err != nil {
		return err
	}
	return nil
}
