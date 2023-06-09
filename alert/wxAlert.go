package alert

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type WxNotify struct {
	webhookUrl      string
	Header          string
	ContentType     string
	MsgTemplate     string
	RecoverTemplate string
}

func NewWxNotify(webhookUrl string) *WxNotify {
	return &WxNotify{
		webhookUrl:  webhookUrl,
		Header:      "application/json;charset=utf-8",
		ContentType: "markdown",
		//MsgTemplate:     MsgTemplate,
		//RecoverTemplate: RecoverTemplate,
	}
}

type AlertMsg struct {
	Alert_name string
	Value      int64
	Time       time.Time
	ExtUrl     string
	Throld     int64
}

func (w *WxNotify) Send(alertMsg *AlertMsg) {

	msg := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"content": "**告警项目**: <font color=\"warning\">" + alertMsg.Alert_name + "</font>\n" +
				"**告警阈值**: " + strconv.Itoa(int(alertMsg.Throld)) + "次/每分钟\n" +
				"**当前值**: " + strconv.Itoa(int(alertMsg.Value)) + "次/每分钟\n" +
				"**告警时间**:" + alertMsg.Time.Format("2006-01-02 15:04:05") + "\n" +
				"[查看错误详情]" + "(" + alertMsg.ExtUrl + ")" + "\n"},
	}

	byteMsg, _ := json.Marshal(msg)

	buffer := bytes.NewBuffer(byteMsg)

	http.Post(w.webhookUrl, w.ContentType, buffer)
}
