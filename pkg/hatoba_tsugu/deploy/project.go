package deploy

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"hatoba_tsugu/pkg/app"
	"hatoba_tsugu/pkg/hatoba_tsugu"
	"hatoba_tsugu/pkg/kubernetes"
	"istio.io/api/networking/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strings"
	"time"
)

func FetchGitProject(git string) (*ProjectCrd, error) {
	response, err := kubernetes.Client.R().SetResult(&ProjectCrd{}).
		Get(kubernetes.HatobaTsuguDeployProjectPath.OnePath(kubernetes.TransferGitDns(git)))
	if err != nil {
		return nil, err
	}

	if response.StatusCode() == http.StatusNotFound {
		return nil, errors.New("not found")
	} else {
		return response.Result().(*ProjectCrd), nil
	}
}

type Project struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Git         *Git      `json:"git"`
	Resource    *Resource `json:"resource"`
}

type ProjectCrd struct {
	v1.TypeMeta   `json:",omitempty"`
	v1.ObjectMeta `json:"metadata,omitempty"`
	Spec          Project `json:"spec,omitempty"`
}

func (p *Project) Launch(env string, image string) error {
	version := time.Now().Format("15-04-05-20060102")
	deployment := p.Deployment(env, version, image)
	response, err := kubernetes.Client.R().SetResult(&appsv1.Deployment{}).
		Get(kubernetes.DeploymentPath.OnePath(deployment.ObjectMeta.Name))
	if err != nil {
		return err
	} else if err = kubernetes.ResponseOk(response); err != nil {
		return err
	}

	if response.StatusCode() == http.StatusNotFound {
		response, err = kubernetes.Client.R().SetBody(deployment).Post(kubernetes.DeploymentPath.MultiPath())
		if err != nil {
			return err
		} else if err = kubernetes.ResponseOk(response); err != nil {
			return err
		}
	} else {
		old := response.Result().(*appsv1.Deployment)
		deployment.ObjectMeta.ResourceVersion = old.ObjectMeta.ResourceVersion
		response, err = kubernetes.Client.R().SetBody(deployment).Put(kubernetes.DeploymentPath.OnePath(deployment.ObjectMeta.Name))
		if err != nil {
			return err
		} else if err = kubernetes.ResponseOk(response); err != nil {
			return err
		}
	}

	service := p.Service(env)
	response, err = kubernetes.Client.R().SetResult(&apiv1.Service{}).
		Get(kubernetes.ServicePath.OnePath(service.ObjectMeta.Name))
	if err != nil {
		return err
	} else if err = kubernetes.ResponseOk(response); err != nil {
		return err
	}

	if response.StatusCode() == http.StatusNotFound {
		response, err = kubernetes.Client.R().SetBody(service).Post(kubernetes.ServicePath.MultiPath())
		if err != nil {
			return err
		} else if err = kubernetes.ResponseOk(response); err != nil {
			return err
		}
	} else {
		old := response.Result().(*apiv1.Service)
		service.ObjectMeta.ResourceVersion = old.ObjectMeta.ResourceVersion
		response, err = kubernetes.Client.R().SetBody(service).Put(kubernetes.ServicePath.OnePath(service.ObjectMeta.Name))
		if err != nil {
			return err
		} else if err = kubernetes.ResponseOk(response); err != nil {
			return err
		}
	}
	return p.AppendDestinationRuleVersion(env, version)
}

func (p *Project) DnsName() string {
	return kubernetes.TransferGitDns(p.Git.Url)
}

func (p *Project) ProjectName(env string) string {
	return fmt.Sprintf("%s-%s", env, p.DnsName())
}

func (p *Project) Labels(env string) map[string]string {
	return map[string]string{
		"env":  env,
		"role": hatoba_tsugu.RoleApp,
		"app":  p.DnsName(),
	}
}

func (p *Project) ConfigMap(env string) (volumes []apiv1.Volume, mounts []apiv1.VolumeMount) {
	volume := apiv1.Volume{}
	volume.ConfigMap = &apiv1.ConfigMapVolumeSource{LocalObjectReference: apiv1.LocalObjectReference{Name: p.ProjectName(env)}}
	volume.Name = hatoba_tsugu.VolumeAppConfig
	volumes = append(volumes, volume)
	for _, config := range p.Resource.Configs {
		if config.Env == env {
			for _, file := range config.Files {
				mount := apiv1.VolumeMount{
					Name:      hatoba_tsugu.VolumeAppConfig,
					MountPath: file.Path,
				}
				paths := strings.Split(file.Path, "/")
				mount.SubPath = paths[len(paths)-1]
				mounts = append(mounts, mount)
			}
		}
	}
	return
}

func (p *Project) ContainerPorts() (ports []apiv1.ContainerPort) {
	for _, port := range p.Resource.Ports {
		por := apiv1.ContainerPort{}
		por.Name = port.Name
		por.ContainerPort = port.TargetPort.IntVal
		por.Protocol = port.Protocol
		ports = append(ports, por)
	}
	return
}

func (p *Project) Service(env string) *apiv1.Service {
	as := &apiv1.Service{
		TypeMeta: v1.TypeMeta{
			Kind:       kubernetes.ServicePath.Kind,
			APIVersion: kubernetes.ServicePath.ApiVersion(),
		},
	}
	labels := p.Labels(env)

	as.ObjectMeta.Name = p.ProjectName(env)
	as.ObjectMeta.Namespace = app.Config.Cd.Namespace
	as.ObjectMeta.Labels = labels
	as.Spec.Selector = labels
	as.Spec.Ports = p.Resource.Ports
	return as
}

func (p *Project) Deployment(env string, version string, image string) *appsv1.Deployment {
	as := &appsv1.Deployment{
		TypeMeta: v1.TypeMeta{
			Kind:       kubernetes.DeploymentPath.Kind,
			APIVersion: kubernetes.DeploymentPath.ApiVersion(),
		},
	}

	// 名字无用
	name := uuid.New().String()
	labels := p.Labels(env)
	labels["version"] = version

	as.ObjectMeta.Name = name
	as.ObjectMeta.Namespace = app.Config.Cd.Namespace
	as.ObjectMeta.Labels = labels
	as.Spec.Selector = &v1.LabelSelector{MatchLabels: labels}
	as.Spec.Template.ObjectMeta.Labels = labels

	configVolumes, configMounts := p.ConfigMap(env)
	as.Spec.Template.Spec.Volumes = configVolumes
	as.Spec.Template.Spec.Containers = []apiv1.Container{
		{
			Name:         name,
			Image:        image,
			VolumeMounts: configMounts,
			Ports:        p.ContainerPorts(),
		},
	}
	return as
}

func (p *Project) AppendDestinationRuleVersion(env string, version string) error {
	dr := &kubernetes.K8sIstioDestinationRule{
		TypeMeta: v1.TypeMeta{
			Kind:       kubernetes.DestinationRulePath.Kind,
			APIVersion: kubernetes.DestinationRulePath.ApiVersion(),
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      p.ProjectName(env),
			Namespace: kubernetes.DestinationRulePath.Ns,
			Labels:    p.Labels(env),
		},
	}
	subset := &v1beta1.Subset{
		Name: version,
		Labels: map[string]string{
			"version": version,
		},
	}
	response, err := kubernetes.Client.R().SetResult(&kubernetes.K8sIstioDestinationRule{}).
		Get(kubernetes.DestinationRulePath.OnePath(p.ProjectName(env)))
	if err != nil {
		return err
	} else if err = kubernetes.ResponseOk(response); err != nil {
		return err
	}

	if response.StatusCode() == http.StatusNotFound {
		dr.Spec.Host = fmt.Sprintf("%s.%s.%s", p.ProjectName(env), app.Config.Cd.Namespace, app.Config.Cd.Domain)
		dr.Spec.Subsets = []*v1beta1.Subset{subset}
		response, err = kubernetes.Client.R().SetBody(dr).Post(kubernetes.DestinationRulePath.MultiPath())
		if err != nil {
			return err
		} else if err = kubernetes.ResponseOk(response); err != nil {
			return err
		}
	} else {
		dr := response.Result().(*kubernetes.K8sIstioDestinationRule)
		dr.Spec.Subsets = append(dr.Spec.Subsets, subset)
		response, err = kubernetes.Client.R().SetBody(dr).
			Put(kubernetes.DestinationRulePath.OnePath(p.ProjectName(env)))
		if err != nil {
			return err
		} else if err = kubernetes.ResponseOk(response); err != nil {
			return err
		}
	}
	return nil
}

type Git struct {
	Url string `json:"url"`
}

type Resource struct {
	Ports      []apiv1.ServicePort `json:"ports"`
	Configs    []*ResourceConfig   `json:"configs"`
	Dockerfile string              `json:"dockerfile"`
}

type ResourceConfig struct {
	Env   string                `json:"env"`
	Files []*ResourceConfigFile `json:"files"`
}

type ResourceConfigFile struct {
	Path        string `json:"path"`
	Config      string `json:"config"`
	Description string `json:"description"`
}
