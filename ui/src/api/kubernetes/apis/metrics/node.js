import maxiloVue from 'maxilo-vue'
let utils = maxiloVue.make("utils")
export default utils.K8sApi(new utils.K8sPath({
    group:   "apis",
    api:     "metrics.k8s.io",
    version: "v1beta1",
    ns:      "",
    kind:    "NodeMetrics",
    plural:  "nodes",
}))