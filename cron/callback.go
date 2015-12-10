package cron

import (
	"fmt"
	"github.com/ZeaLoVe/alarm/api"
	"github.com/ZeaLoVe/alarm/redis"
	"github.com/open-falcon/common/model"
	"github.com/toolkits/net/httplib"
	"strings"
	"time"
)

func HandleCallback(event *model.Event, action *api.Action) {

	teams := action.Uic
	mails := []string{}
	ims := []string{}

	if teams != "" {
		_, mails, ims = api.ParseTeams(teams)
		imsmsContent := GenerateIMSmsContent(event)
		mailContent := GenerateMailContent(event)
		if action.BeforeCallbackSms == 1 {
			redis.WriteIMSms(ims, imsmsContent)
		}

		if action.BeforeCallbackMail == 1 {
			redis.WriteMail(mails, imsmsContent, mailContent)
		}
	}

	message := Callback(event, action)

	if teams != "" {
		if action.AfterCallbackSms == 1 {
			redis.WriteIMSms(ims, message)
		}

		if action.AfterCallbackMail == 1 {
			redis.WriteMail(mails, message, message)
		}
	}

}

func Callback(event *model.Event, action *api.Action) string {
	if action.Url == "" {
		return "callback url is blank"
	}

	L := make([]string, 0)
	if len(event.PushedTags) > 0 {
		for k, v := range event.PushedTags {
			L = append(L, fmt.Sprintf("%s:%s", k, v))
		}
	}

	tags := ""
	if len(L) > 0 {
		tags = strings.Join(L, ",")
	}

	req := httplib.Get(action.Url).SetTimeout(3*time.Second, 20*time.Second)

	req.Param("endpoint", event.Endpoint)
	req.Param("metric", event.Metric())
	req.Param("status", event.Status)
	req.Param("step", fmt.Sprintf("%d", event.CurrentStep))
	req.Param("priority", fmt.Sprintf("%d", event.Priority()))
	req.Param("time", event.FormattedTime())
	req.Param("tpl_id", fmt.Sprintf("%d", event.TplId()))
	req.Param("exp_id", fmt.Sprintf("%d", event.ExpressionId()))
	req.Param("stra_id", fmt.Sprintf("%d", event.StrategyId()))
	req.Param("tags", tags)

	resp, e := req.String()

	success := "success"
	if e != nil {
		success = fmt.Sprintf("fail:%s", e.Error())
	}
	message := fmt.Sprintf("curl %s %s. resp: %s", action.Url, success, resp)

	return message
}
