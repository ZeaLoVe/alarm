package cron

import (
	"fmt"
	"github.com/ZeaLoVe/alarm/g"
	"github.com/open-falcon/common/model"
	"github.com/open-falcon/common/utils"
)

func BuildCommonSMSContent(event *model.Event) string {
	return fmt.Sprintf(
		"[P%d][%s][Endpoint:%s][][%s %s %s %s %s%s%s][n%d %s]",
		event.Priority(),
		event.Status,
		event.Endpoint,
		event.Note(),
		event.Func(),
		event.Metric(),
		utils.SortedTags(event.PushedTags),
		utils.ReadableFloat(event.LeftValue),
		event.Operator(),
		utils.ReadableFloat(event.RightValue()),
		event.CurrentStep,
		event.FormattedTime(),
	)
}

func BuildCommonIMSmsContent(event *model.Event) string {
	return fmt.Sprintf(
		"[P%d][%s][Endpoint:%s][][%s %s %s %s %s%s%s][累计出现 %d 次 时间 %s]",
		event.Priority(),
		event.Status,
		event.Endpoint,
		event.Note(),
		event.Func(),
		event.Metric(),
		utils.SortedTags(event.PushedTags),
		utils.ReadableFloat(event.LeftValue),
		event.Operator(),
		utils.ReadableFloat(event.RightValue()),
		event.CurrentStep,
		event.FormattedTime(),
	)
}

func BuildCommonPhoneContent(event *model.Event) string {
	return fmt.Sprintf(
		"Endpoint:%s,你所设置的指标%s出现故障:%s,该指标值为%s,报警值为%s.该报警累计出现%d次",
		event.Endpoint,
		event.Metric(),
		event.Note(),
		utils.ReadableFloat(event.LeftValue),
		utils.ReadableFloat(event.RightValue()),
		event.CurrentStep,
	)
}

func BuildCommonMailContent(event *model.Event) string {
	link := g.Link(event)
	return fmt.Sprintf(
		"%s\r\nP%d\r\nEndpoint:%s\r\nMetric:%s\r\nTags:%s\r\n%s: %s%s%s\r\nNote:%s\r\nMax:%d, Current:%d\r\nTimestamp:%s\r\n%s\r\n",
		event.Status,
		event.Priority(),
		event.Endpoint,
		event.Metric(),
		utils.SortedTags(event.PushedTags),
		event.Func(),
		utils.ReadableFloat(event.LeftValue),
		event.Operator(),
		utils.ReadableFloat(event.RightValue()),
		event.Note(),
		event.MaxStep(),
		event.CurrentStep,
		event.FormattedTime(),
		link,
	)
}

func GenerateSmsContent(event *model.Event) string {
	return BuildCommonSMSContent(event)
}

func GenerateMailContent(event *model.Event) string {
	return BuildCommonMailContent(event)
}

func GenerateIMSmsContent(event *model.Event) string {
	return BuildCommonIMSmsContent(event)
}

func GeneratePhoneContent(event *model.Event) string {
	return BuildCommonPhoneContent(event)
}
