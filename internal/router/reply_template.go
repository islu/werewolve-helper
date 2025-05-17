package router

import "github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"

func ModeSettingTemplateV2(liffID string) messaging_api.MessageInterface {
	return &messaging_api.TemplateMessage{
		AltText: "Setting role alt text",
		Template: &messaging_api.ButtonsTemplate{
			Title: "開設房間",
			Text:  "請點擊開始設定",
			Actions: []messaging_api.ActionInterface{
				messaging_api.UriAction{Label: "開始設定", Uri: "https://liff.line.me/" + liffID},
			},
		},
	}
}
