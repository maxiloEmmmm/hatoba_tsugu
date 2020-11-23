package main

import (
	"github.com/gin-gonic/gin"
	"github.com/maxiloEmmmm/go-web/contact"
	"hatoba_tsugu/api/deploy"
	"hatoba_tsugu/pkg/kubernetes"
)

func InitRoute() *gin.Engine {
	engine := gin.New()
	engine.Use(contact.GinCors())
	engine.Any("/cloud-api/*path", kubernetes.UiProxy())

	engine.POST("/launch", contact.GinHelpHandle(deploy.ProjectLaunch))
	engine.GET("/dockerfile", contact.GinHelpHandle(deploy.ProjectBuildConf))
	return engine
}
