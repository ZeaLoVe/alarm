package cron

import (
	"github.com/ZeaLoVe/alarm/api"
	"github.com/ZeaLoVe/alarm/redis"
	"github.com/open-falcon/common/model"
)

func consume(event *model.Event, isHigh bool) {
	actionId := event.ActionId()
	if actionId <= 0 {
		return
	}

	action := api.GetAction(actionId)
	if action == nil {
		return
	}

	if action.Callback == 1 {
		HandleCallback(event, action)
		return
	}

	consumeEvents(event, action)

}

// 处理事件（去掉了优先级区分）
func consumeEvents(event *model.Event, action *api.Action) {
	if action.Uic == "" {
		return
	}
	//获取联系方式
	phones, mails, ims := api.ParseTeams(action.Uic)

	smsContent := GenerateSmsContent(event)
	mailContent := GenerateMailContent(event)
	imContent := GenerateIMSmsContent(event)
	phoneContent := GeneratePhoneContent(event)

	// level 0 will have phone,im,mail
	// level 1-2 will have im,mail
	// level 3-6 will only have im

	if event.Priority() == 0 {
		redis.WritePhone(phones, phoneContent)
	}

	if event.Priority() < 3 {
		redis.WriteMail(mails, smsContent, mailContent)
	}

	redis.WriteIMSms(ims, imContent)
}
