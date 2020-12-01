package app

import (
	"github.com/maxiloEmmmm/go-web/contact"
	"log"
)

type config struct {
	Kubernetes KubernetesConfig
	Cd         Cd
	Istio      Istio
}

type KubernetesConfig struct {
	ApiServer string `yaml:"api_server"`
	Token     string `yaml:"token"`
}

type Cd struct {
	Namespace string
	Domain    string
}

type Istio struct {
	Kiali string `json:"kiali"`
}

var Config *config

func initConfig() {
	Config = &config{}
	ok := contact.ConfigFile("app.yaml", Config)

	if !ok {
		log.Printf("可 复制app.yaml.example创建app.yaml\n")
	}
}
