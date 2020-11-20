package main

import (
	"github.com/gin-gonic/gin"
	"github.com/maxiloEmmmm/go-web/contact"
	"hatoba_tsugu/pkg/kubernetes"
)

func InitRoute() *gin.Engine {
	engine := gin.New()
	engine.Use(contact.GinCors())
	engine.Any("/cloud-api/*path", kubernetes.UiProxy())
	return engine
}
