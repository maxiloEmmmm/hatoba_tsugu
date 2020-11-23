package kubernetes

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"hatoba_tsugu/pkg/app"
	"istio.io/api/networking/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

var (
	Client *resty.Client

	HatobaTsuguDeployProjectPath K8Path
	DeploymentPath K8Path
	ServicePath K8Path
	DestinationRulePath K8Path
)

type K8sIstioDestinationRule struct {
	v1.TypeMeta `json:",omitempty"`
	v1.ObjectMeta `json:"metadata,omitempty"`
	Spec v1beta1.DestinationRule `json:"spec"`
}

func Init() {
	Client = resty.New()
	Client.SetAuthToken(getToken())
	Client.HostURL = app.Config.Kubernetes.ApiServer

	HatobaTsuguDeployProjectPath = K8Path{
		Group: "apis",
		Api: "deploy.hatobatsugu.gsc",
		Version: "v1",
		Ns: app.Config.Cd.Namespace,
		Kind: "projects",
	}
	DeploymentPath = K8Path{
		Group: "apis",
		Api: "apps",
		Version: "v1",
		Ns: app.Config.Cd.Namespace,
		Kind: "deployments",
	}
	ServicePath = K8Path{
		Group: "api",
		Api: "",
		Version: "v1",
		Ns: app.Config.Cd.Namespace,
		Kind: "services",
	}
	DestinationRulePath = K8Path{
		Group: "apis",
		Api: "networking.istio.io",
		Version: "v1beta1",
		Ns: app.Config.Cd.Namespace,
		Kind: "destinationrules",
	}
}

type K8Path struct {
	Group string
	Api string
	Version string
	Kind string
	Ns string
}

func (k8p K8Path) MultiPath() string {
	return fmt.Sprintf("/%s/%s/%s/namespaces/%s/%s", k8p.Group, k8p.Api, k8p.Version, k8p.Ns, k8p.Kind)
}

func (k8p K8Path) OnePath(name string) string {
	return fmt.Sprintf("/%s/%s", k8p.MultiPath(), name)
}

func (k8p K8Path) ApiVersion() string {
	return fmt.Sprintf("%s/%s", k8p.Api, k8p.Kind)
}

func FullUpdateOrCreate(point K8Path, name string, data interface{}) error {
	response, err := Client.R().
		Get(point.OnePath(name))
	if err != nil {
		return err
	}

	if response.StatusCode() == http.StatusNotFound {
		response, err = Client.R().SetBody(data).Post(point.MultiPath())
		if err != nil {
			return err
		}
	}else {
		response, err = Client.R().SetBody(data).Put(point.OnePath(name))
		if err != nil {
			return err
		}
	}

	return nil
}
