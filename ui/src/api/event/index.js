import maxiloVue from 'maxilo-vue'

export default {
    channel(){
        return maxiloVue.make("http").get("/channel")
    },
    filter(){
        return maxiloVue.make("http").get("/notification/filter")
    },
    refreshFilter(){
        return maxiloVue.make("http").get("/notification/filter/refresh")
    }
}