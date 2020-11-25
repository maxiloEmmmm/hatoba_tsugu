import maxiloVue from 'maxilo-vue'
let config = maxiloVue.make("config")
config.add("baseURL", process.env.VUE_APP_BASEURL ? process.env.VUE_APP_BASEURL : "")

let base = {
    cd_ns: "apps",
    cd_domain: "svc.cluster.local",
    cd_endpoint_ns: "endpoint"
}
config.merge(base)

export default base