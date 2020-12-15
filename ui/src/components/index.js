import maxiloVue from "maxilo-vue"

import istioVS from "./istioVS"
maxiloVue.vue.component(istioVS.name, istioVS)

import projectPort from "./project-port"
maxiloVue.vue.component(projectPort.name, projectPort)

import projectConfig from "./project-config"
maxiloVue.vue.component(projectConfig.name, projectConfig)

import projectPrometheus from "./project-prometheus"
maxiloVue.vue.component(projectPrometheus.name, projectPrometheus)

import notificationFilter from "./notification-filter"
maxiloVue.vue.component(notificationFilter.name, notificationFilter)

import log from "./log"
maxiloVue.vue.component(log.name, log)