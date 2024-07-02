package internal

import "github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"

func ModeSettingTemplate() messaging_api.MessageInterface {
	return &messaging_api.TemplateMessage{
		AltText: "Setting game mode image carousel alt text",
		Template: &messaging_api.ImageCarouselTemplate{
			Columns: []messaging_api.ImageCarouselColumn{
				{
					ImageUrl: "https://raw.githubusercontent.com/islu/Werewolve-Helper/main/img/9_classic_mode.png",
					Action:   messaging_api.PostbackAction{Label: "9人標準配置", Data: "9人標準配置"},
				},
				{
					ImageUrl: "https://raw.githubusercontent.com/islu/Werewolve-Helper/main/img/all-seeing-eye.png",
					Action:   messaging_api.PostbackAction{Label: "自訂配置", Data: "自訂配置"},
				},
			},
		},
	}
}

func CustomModeTemplate() messaging_api.MessageInterface {
	return &messaging_api.TemplateMessage{
		AltText: "Custom mode carousel alt text",
		Template: &messaging_api.CarouselTemplate{
			Columns: []messaging_api.CarouselColumn{
				{
					Title: "自訂配置",
					Text:  "設定人數",
					Actions: []messaging_api.ActionInterface{
						messaging_api.PostbackAction{Label: "設定預言家", Data: "設定預言家"},
						messaging_api.PostbackAction{Label: "設定女巫", Data: "設定女巫"},
						messaging_api.PostbackAction{Label: "設定獵人", Data: "設定獵人"},
					},
				},
				{
					Title: "自訂配置",
					Text:  "設定人數",
					Actions: []messaging_api.ActionInterface{
						messaging_api.PostbackAction{Label: "設定平民", Data: "設定平民"},
						messaging_api.PostbackAction{Label: "設定狼人", Data: "設定狼人"},
						// messaging_api.PostbackAction{Label: "設定白癡", Data: "設定白癡"},
					},
				},
			},
		},
	}
}

func QuickReplyButtons() messaging_api.MessageInterface {
	return &messaging_api.TextMessage{
		Text: "設定人數",
		QuickReply: &messaging_api.QuickReply{
			Items: []messaging_api.QuickReplyItem{
				{Action: &messaging_api.MessageAction{Label: "1", Text: "1"}},
				{Action: &messaging_api.MessageAction{Label: "2", Text: "2"}},
				{Action: &messaging_api.MessageAction{Label: "3", Text: "3"}},
				{Action: &messaging_api.MessageAction{Label: "4", Text: "4"}},
				{Action: &messaging_api.MessageAction{Label: "5", Text: "5"}},
				{Action: &messaging_api.MessageAction{Label: "6", Text: "6"}},
				{Action: &messaging_api.MessageAction{Label: "7", Text: "7"}},
				{Action: &messaging_api.MessageAction{Label: "8", Text: "8"}},
				{Action: &messaging_api.MessageAction{Label: "9", Text: "9"}},
				{Action: &messaging_api.MessageAction{Label: "10", Text: "10"}},
				{Action: &messaging_api.MessageAction{Label: "11", Text: "11"}},
				{Action: &messaging_api.MessageAction{Label: "12", Text: "12"}},
				{Action: &messaging_api.MessageAction{Label: "13", Text: "13"}},
			},
		},
	}
}
