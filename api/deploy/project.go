package deploy

import (
	"github.com/maxiloEmmmm/go-web/contact"
	"hatoba_tsugu/api/types"
	"hatoba_tsugu/pkg/hatoba_tsugu"
	"hatoba_tsugu/pkg/hatoba_tsugu/deploy"
)

func ProjectLaunch(help *contact.GinHelp) {
	body := &types.Launch{}
	help.InValidBind(body)

	body.Env = hatoba_tsugu.TransferEnv(body.Env)
	projectCrd, err := deploy.FetchGitProject(body.Git)
	help.AssetsInValid("fetch", err)
	if err = projectCrd.Spec.Launch(body.Env, body.Image); err != nil {
		projectCrd.Spec.LaunchFailEvent(err)
	} else {
		projectCrd.Spec.LaunchSuccessFailEvent()
	}
	help.AssetsInValid("launch", err)
	help.Resource(nil)
}

func ProjectBuildConf(help *contact.GinHelp) {
	projectCrd, err := deploy.FetchGitProject(help.DefaultQuery("git", ""))
	help.AssetsInValid("fetch", err)
	help.String(200, "%s", projectCrd.Spec.Resource.Dockerfile)
}
