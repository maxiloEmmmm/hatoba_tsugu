package main

import (
	"hatoba_tsugu/pkg/app"
	"hatoba_tsugu/pkg/hatoba_tsugu"
	"hatoba_tsugu/pkg/kubernetes"
)

func main() {
	// todo： use context
	stop := make(chan struct{})
	defer close(stop)

	app.Init()
	kubernetes.Init()
	hatoba_tsugu.Init(stop)

	InitRoute().Run(":8000")
}
