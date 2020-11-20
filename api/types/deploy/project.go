package deploy

import (
	apiv1 "k8s.io/api/core/v1"
)

type Project struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Git         *Git      `json:"git"`
	Resource    *Resource `json:"resource"`
}

type Git struct {
	Url string `json:"url"`
}

type Resource struct {
	Ports   []*apiv1.ServicePort `json:"ports"`
	Configs []*ResourceConfig    `json:"configs"`
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
