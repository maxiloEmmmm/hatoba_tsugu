import maxiloVue from 'maxilo-vue'
import config from "@/config"
let utils = maxiloVue.make("utils")
export default utils.K8sApi(new utils.K8sPath({
    group:   "apis",
    api:     "metrics.k8s.io",
    version: "v1beta1",
    ns:      config.cd_ns,
    kind:    "PodMetrics",
    plural:  "pods",
}))