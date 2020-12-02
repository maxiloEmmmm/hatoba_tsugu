package channel

import (
	"github.com/go-resty/resty/v2"
	"github.com/maxiloEmmmm/go-web/contact"
)

type WebhookChannel struct {
	Config
}

func (ddc *WebhookChannel) Send(msg string) bool {
	_, err := resty.New().R().
		SetHeader("Content-Type", "application/json;charset=utf-8").
		SetBody(contact.H{
			"msg": msg,
		}).Post(ddc.Config["url"])
	if err != nil {
		contact.Warning.Log("evt.channel.webhook", err.Error())
		return false
	}
	return false
}
