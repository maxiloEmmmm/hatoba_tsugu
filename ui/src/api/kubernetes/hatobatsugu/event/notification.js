import maxiloVue from 'maxilo-vue'
import config from "@/config"
let utils = maxiloVue.make("utils")
export default utils.K8sApi(new utils.K8sPath({
    group:   "apis",
    api:     "events.hatobatsugu.gsc",
    version: "v1",
    ns:      config.cd_ns,
    kind:    "Notification",
    plural:  "notifications",
}))