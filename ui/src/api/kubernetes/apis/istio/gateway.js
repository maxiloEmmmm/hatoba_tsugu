import maxiloVue from 'maxilo-vue'
import config from "@/config"
let utils = maxiloVue.make("utils")
export default utils.K8sApi(new utils.K8sPath({
    group:   "apis",
    api:     "networking.istio.io",
    version: "v1beta1",
    ns:      config.cd_endpoint_ns,
    kind:    "Gateway",
    plural:  "gateways",
}))