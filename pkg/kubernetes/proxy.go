package kubernetes

import (
	"crypto/tls"
	"github.com/gin-gonic/gin"
	go_tool "github.com/maxiloEmmmm/go-tool"
	"github.com/maxiloEmmmm/go-web/contact"
	"hatoba_tsugu/pkg/app"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func UiProxy() gin.HandlerFunc {
	u, err := url.Parse(app.Config.Kubernetes.ApiServer)
	if err != nil {
		log.Fatalf("kubernetes apiserver parse err, %s", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	token := getToken()
	return contact.GinHelpHandle(func(c *contact.GinHelp) {
		uri := &struct {
			Path string `uri:"path"`
		}{}
		c.InValidBindUri(uri)
		c.Request.Header.Set("Authorization", go_tool.StringJoin("Bearer ", token))
		c.Request.URL.Path = uri.Path
		proxy.ServeHTTP(c.Writer, c.Request)
	})
}

func getToken() string {
	token := app.Config.Kubernetes.Token
	if token == "" {
		data, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
		if err != nil {
			log.Fatalf("kubernetes get container token err, %s", err)
		}

		token = string(data)
	}
	return token
}
