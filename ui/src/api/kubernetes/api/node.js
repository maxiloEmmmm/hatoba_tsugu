import maxiloVue from 'maxilo-vue'
let utils = maxiloVue.make("utils")
export default utils.K8sApi(new utils.K8sPath({
    group:   "api",
    api:     "",
    version: "v1",
    kind:    "Node",
    plural:  "nodes",
}))