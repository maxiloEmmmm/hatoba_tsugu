import maxiloVue from 'maxilo-vue'
import tool from 'antd-vue-tool'
import axios from "axios"
let config = maxiloVue.make("config")
maxiloVue.register({
    register: function(){},
    boot: function(app){
        let http = axios.create({
            baseURL: `${config.baseURL ? config.baseURL : `${window.location.protocol}//${window.location.host}`}/cloud-api`
        })

        Object.defineProperty(app.vue.prototype, '$kb', {
            get: () => {
                return http;
            }
        })
        tool.http.ext.cloud = {
            engine: http,
            errorMsgAdapter: r => {
                return `「${r.response.data.code}」 ${r.response.data.message}`
            }
        }
    }
})

maxiloVue.register({
    register: function(){},
    boot: function(app){
        tool.http.engine = app.make("http")
        tool.http.errorMsgAdapter = r => {
            return `「${r.response.data.code}」 ${r.response.data.message}`
        }
    }
})
