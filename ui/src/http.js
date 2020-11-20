import maxiloVue from 'maxilo-vue'
import tool from 'antd-vue-tool'
import axios from "axios"
import {merge} from "lodash"
let config = maxiloVue.make("config")
config.add("baseURL", process.env.VUE_APP_BASEURL ? process.env.VUE_APP_BASEURL : "http://api")
maxiloVue.register({
    register: function(){},
    boot: function(app){
        let http = axios.create({
            baseURL: `${config.baseURL}/cloud-api`
        })

        http.fullUpdateOrCreate = async (obj, data) => {
            try {
                let response = await http.get(`${obj}/${data.metadata.name}`)
                await http.put(`${obj}/${data.metadata.name}`, merge(response.data, data))
            }catch(e) {
                if(e.response.data.code == 404) {
                    try {
                        await http.post(obj, data)
                    }catch(e) {
                        throw e
                    }
                }else {
                    throw e
                }
            }
        }
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
