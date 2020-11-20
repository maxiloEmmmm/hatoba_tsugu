import maxiloVue from "maxilo-vue"

import istioVS from "./istioVS"
maxiloVue.vue.component(istioVS.name, istioVS)

import projectPort from "./project-port"
maxiloVue.vue.component(projectPort.name, projectPort)

import projectConfig from "./project-config"
maxiloVue.vue.component(projectConfig.name, projectConfig)