package kubernetes

import (
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	go_tool "github.com/maxiloEmmmm/go-tool"
	"hatoba_tsugu/pkg/app"
	"istio.io/api/networking/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"net/http"
)

var (
	Client                       *resty.Client
	KubeClient                   *kubernetes.Clientset
	HatobaTsuguDeployProjectPath K8Path

	HatobaTsuguEventNotificationPath K8Path
	DeploymentPath                   K8Path
	ServicePath                      K8Path
	EventPath                        K8Path
	DestinationRulePath              K8Path
)

type K8sIstioDestinationRule struct {
	v1.TypeMeta   `json:",omitempty"`
	v1.ObjectMeta `json:"metadata,omitempty"`
	Spec          v1beta1.DestinationRule `json:"spec"`
}

func Init() {
	var err error
	KubeClient, err = kubernetes.NewForConfig(&rest.Config{
		Host:        app.Config.Kubernetes.ApiServer,
		BearerToken: getToken(),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	})

	if err != nil {
		log.Fatal(fmt.Sprintf("kube client: %s", err))
	}

	Client = resty.New()
	Client.SetAuthToken(getToken())
	Client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	Client.HostURL = app.Config.Kubernetes.ApiServer

	HatobaTsuguDeployProjectPath = K8Path{
		Group:   "apis",
		Api:     "deploy.hatobatsugu.gsc",
		Version: "v1",
		Ns:      app.Config.Cd.Namespace,
		Kind:    "Project",
		Plural:  "projects",
	}
	HatobaTsuguEventNotificationPath = K8Path{
		Group:   "apis",
		Api:     "events.hatobatsugu.gsc",
		Version: "v1",
		Ns:      app.Config.Cd.Namespace,
		Kind:    "Notification",
		Plural:  "notifications",
	}
	DeploymentPath = K8Path{
		Group:   "apis",
		Api:     "apps",
		Version: "v1",
		Ns:      app.Config.Cd.Namespace,
		Kind:    "Deployment",
		Plural:  "deployments",
	}
	ServicePath = K8Path{
		Group:   "api",
		Api:     "",
		Version: "v1",
		Ns:      app.Config.Cd.Namespace,
		Kind:    "Service",
		Plural:  "services",
	}
	EventPath = K8Path{
		Group:   "api",
		Api:     "",
		Version: "v1",
		Ns:      app.Config.Cd.Namespace,
		Kind:    "Event",
		Plural:  "events",
	}
	DestinationRulePath = K8Path{
		Group:   "apis",
		Api:     "networking.istio.io",
		Version: "v1beta1",
		Ns:      app.Config.Cd.Namespace,
		Kind:    "DestinationRule",
		Plural:  "destinationrules",
	}
}

type K8Path struct {
	Group   string
	Api     string
	Version string
	Kind    string
	Plural  string
	Ns      string
}

func (k8p K8Path) MultiPath() string {
	return fmt.Sprintf("/%s%s/%s/namespaces/%s/%s", k8p.Group, go_tool.AssetsReturn(k8p.Api == "", "", go_tool.StringJoin("/", k8p.Api)), k8p.Version, k8p.Ns, k8p.Plural)
}

func (k8p K8Path) OnePath(name string) string {
	return fmt.Sprintf("%s/%s", k8p.MultiPath(), name)
}

func (k8p K8Path) ApiVersion() string {
	return fmt.Sprintf("%s%s", go_tool.AssetsReturn(k8p.Api == "", "", go_tool.StringJoin(k8p.Api, "/")), k8p.Version)
}

type InValidMessage struct {
	Message string
	Code    int
}
