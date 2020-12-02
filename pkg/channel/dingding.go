package channel

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/go-resty/resty/v2"
	go_tool "github.com/maxiloEmmmm/go-tool"
	"github.com/maxiloEmmmm/go-web/contact"
	"net/url"
	"strconv"
	"time"
)

type DingDingChannel struct {
	Config
}

func (ddc *DingDingChannel) Send(msg string) bool {
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	secret := ddc.Config["secret"]
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(go_tool.StringJoin(timestamp, "\n", secret)))

	_, err := resty.New().R().
		SetHeader("Content-Type", "application/json;charset=utf-8").
		SetBody(contact.H{
			"msgtype": "markdown",
			"markdown": contact.H{
				"title": "evt",
				"text":  msg,
			},
		}).
		SetQueryParams(map[string]string{
			"access_token": ddc.Config["token"],
			"timestamp":    timestamp,
			"sign":         url.QueryEscape(base64.StdEncoding.EncodeToString(h.Sum(nil))),
		}).Post("https://oapi.dingtalk.com/robot/send")
	if err != nil {
		contact.Warning.Log("evt.channel.dingding", err.Error())
		return false
	}
	return true
}
