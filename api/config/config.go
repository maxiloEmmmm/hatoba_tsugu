package config

import (
	"github.com/maxiloEmmmm/go-web/contact"
	"hatoba_tsugu/pkg/app"
)

func Config(help *contact.GinHelp) {
	help.Resource(contact.H{
		"istio": app.Config.Istio,
	})
}
