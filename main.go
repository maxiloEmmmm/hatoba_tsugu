package main

import (
	"github.com/maxiloEmmmm/go-web/contact"
	"hatoba_tsugu/pkg/app"
	"hatoba_tsugu/pkg/hatoba_tsugu"
	"hatoba_tsugu/pkg/kubernetes"
)

func main() {
	// todo： use context
	stop := make(chan struct{})
	defer close(stop)

	contact.Init()
	defer contact.Close()
	app.Init()
	kubernetes.Init()
	hatoba_tsugu.Init(stop)

	InitRoute().Run(":8000")
}
