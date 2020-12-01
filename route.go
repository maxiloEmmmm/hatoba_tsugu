package main

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/maxiloEmmmm/go-web/contact"
	"hatoba_tsugu/api/config"
	"hatoba_tsugu/api/deploy"
	"hatoba_tsugu/pkg/kubernetes"
	"net/http"
)

func InitRoute() *gin.Engine {
	engine := gin.New()
	engine.Use(contact.GinCors())

	if contact.Config.App.Mode == gin.ReleaseMode {
		engine.Use(gzip.Gzip(gzip.DefaultCompression))
		engine.StaticFS("/ui", http.Dir("./front"))
	}

	engine.Any("/cloud-api/*path", kubernetes.UiProxy())
	engine.Any("/cloud-ws/*path", kubernetes.WsProxy())

	engine.POST("/launch", contact.GinHelpHandle(deploy.ProjectLaunch))
	engine.GET("/dockerfile", contact.GinHelpHandle(deploy.ProjectBuildConf))
	engine.GET("/config", contact.GinHelpHandle(config.Config))
	return engine
}
