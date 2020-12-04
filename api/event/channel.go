package event

import (
	go_tool "github.com/maxiloEmmmm/go-tool"
	"github.com/maxiloEmmmm/go-web/contact"
	"hatoba_tsugu/pkg/app"
	"hatoba_tsugu/pkg/hatoba_tsugu/event"
)

func Channel(help *contact.GinHelp) {
	help.Resource(go_tool.ArrayMap(app.Config.Channel, func(d interface{}) interface{} {
		c := d.(app.Channel)
		return contact.H{
			"label": go_tool.StringJoin("[", c.Type, "] ", c.Name),
			"value": c.Name,
		}
	}))
}

func Filter(help *contact.GinHelp) {
	help.Resource([]map[string]string{
		{"label": event.EqFilterType, "value": event.EqFilterType},
		{"label": event.InFilterType, "value": event.InFilterType},
		{"label": event.IgnoreFilterType, "value": event.IgnoreFilterType},
	})
}

func RefreshFilter(help *contact.GinHelp) {
	event.FetchEventNotification()
	help.Resource(nil)
}
