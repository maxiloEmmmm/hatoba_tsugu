import kubernetes from "./kubernetes"
import config from "./config"
import event from "./event"
import maxiloVue from "maxilo-vue"
maxiloVue.register({
    register: function(){},
    boot: function(app){
        Object.defineProperty(app.vue.prototype, '$api', {
            get: () => {
                return {
                    kubernetes,
                    config,
                    event
                }
            }
        })
    }
})