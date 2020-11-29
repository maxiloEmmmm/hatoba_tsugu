import maxiloVue from 'maxilo-vue'
import config from "@/config"
let utils = maxiloVue.make("utils")
export default utils.K8sApi(new utils.K8sPath({
    group:   "api",
    api:     "",
    version: "v1",
    ns:      config.cd_ns,
    kind:    "Pod",
    plural:  "pods",
}))