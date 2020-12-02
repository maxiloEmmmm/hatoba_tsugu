package channel

import "hatoba_tsugu/pkg/app"

const (
	DingDingEngine = "dingding"
	WebhhookEngine = "webhook"
)

type Channel interface {
	Send(msg string) bool
}

func NewChannel(name string) Channel {
	c := getConfigByName(name)
	if c == nil {
		return nil
	}

	switch c.Type {
	case DingDingEngine:
		return &DingDingChannel{c.Config}
	case WebhhookEngine:
		return &WebhookChannel{c.Config}
	}
	return nil
}

func getConfigByName(name string) *app.Channel {
	for _, c := range app.Config.Channel {
		if c.Name == name {
			return &c
		}
	}
	return nil
}

type Config map[string]string
