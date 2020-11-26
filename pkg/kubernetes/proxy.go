package kubernetes

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	go_tool "github.com/maxiloEmmmm/go-tool"
	"github.com/maxiloEmmmm/go-web/contact"
	"hatoba_tsugu/pkg/app"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsProxy() gin.HandlerFunc {
	return contact.GinHelpHandle(func(c *contact.GinHelp) {
		uri := &struct {
			Path string `uri:"path"`
		}{}
		c.InValidBindUri(uri)

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(fmt.Sprintf("%s: %s", "upgrade", err))
			return
		}
		defer conn.Close()

		prefixReg, _ := regexp.Compile("^https?://")
		u := url.URL{Scheme: go_tool.AssetsReturn(strings.HasPrefix(app.Config.Kubernetes.ApiServer, "https"), "wss", "ws").(string), Host: prefixReg.ReplaceAllString(app.Config.Kubernetes.ApiServer, ""), Path: uri.Path, RawQuery: c.Request.URL.RawQuery}
		header := make(http.Header)
		headerWithAuthorization(&header)
		for k, v := range c.Request.Header {
			// https://github.com/gorilla/websocket/blob/c3dd95aea9779669bb3daafbd84ee0530c8ce1c1/client.go#L213
			if !go_tool.InArray([]string{"Upgrade", "Connection", "Sec-Websocket-Key", "Sec-Websocket-Version", "Sec-Websocket-Extensions", "Sec-Websocket-Protocol"}, k) {
				header.Set(k, v[0])
			}
		}

		dialer := websocket.Dialer{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		remoteConn, _, err := dialer.Dial(u.String(), header)
		if err != nil {
			log.Println(fmt.Sprintf("%s: %s", "dial", err))
			return
		}
		defer remoteConn.Close()

		errorChan := make(chan error, 1)
		go func() {
			for {
				mt, message, err := conn.ReadMessage()
				if err != nil {
					errorChan <- errors.New(go_tool.StringJoin("client:read", err.Error()))
					break
				}
				err = remoteConn.WriteMessage(mt, message)
				if err != nil {
					errorChan <- errors.New(go_tool.StringJoin("remote.write", err.Error()))
					break
				}
			}
		}()

		go func() {
			for {
				mt, msg, err := remoteConn.ReadMessage()
				if err != nil {
					errorChan <- errors.New(go_tool.StringJoin("remote.read", err.Error()))
					break
				}
				err = conn.WriteMessage(mt, msg)
				if err != nil {
					errorChan <- errors.New(go_tool.StringJoin("client:write", err.Error()))
					break
				}
			}
		}()

		err = <-errorChan
	})
}

func UiProxy() gin.HandlerFunc {
	u, err := url.Parse(app.Config.Kubernetes.ApiServer)
	if err != nil {
		log.Fatalf("kubernetes apiserver parse err, %s", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return contact.GinHelpHandle(func(c *contact.GinHelp) {
		uri := &struct {
			Path string `uri:"path"`
		}{}
		c.InValidBindUri(uri)
		headerWithAuthorization(&c.Request.Header)
		c.Request.URL.Path = uri.Path
		proxy.ServeHTTP(c.Writer, c.Request)
	})
}

func headerWithAuthorization(header *http.Header) {
	token := getToken()
	header.Set("Authorization", go_tool.StringJoin("Bearer ", token))
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
