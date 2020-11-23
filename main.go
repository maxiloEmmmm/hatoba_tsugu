package main

import (
	"hatoba_tsugu/pkg/app"
	"hatoba_tsugu/pkg/kubernetes"
)

func main() {
	app.Init()
	kubernetes.Init()
	InitRoute().Run(":8000")
}
