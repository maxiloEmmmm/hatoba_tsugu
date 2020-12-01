import maxiloVue from 'maxilo-vue'

export default {
    get(){
        return maxiloVue.make("http").get("/config")
    }
}