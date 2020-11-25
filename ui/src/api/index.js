import kubernetes from "./kubernetes"
import maxiloVue from "maxilo-vue"
maxiloVue.register({
    register: function(){},
    boot: function(app){
        Object.defineProperty(app.vue.prototype, '$api', {
            get: () => {
                return {
                    kubernetes
                }
            }
        })
    }
})