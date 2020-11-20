package main

import "hatoba_tsugu/pkg/app"

func main() {
	app.Init()
	InitRoute().Run(":8000")
}
