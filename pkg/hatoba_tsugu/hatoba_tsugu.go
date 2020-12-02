package hatoba_tsugu

import "hatoba_tsugu/pkg/hatoba_tsugu/event"

func Init(stop chan struct{}) {
	event.Init(stop)
}
