package hatoba_tsugu

import go_tool "github.com/maxiloEmmmm/go-tool"

const (
	ProdEnv = "prod"
	DevEnv = "dev"
)

var Envs = []string{ProdEnv, DevEnv}

func TransferEnv(env string) string {
	if go_tool.InArray(Envs, env) {
		return env
	}else {
		return DevEnv
	}
}