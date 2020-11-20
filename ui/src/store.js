const state = {
    token: ""
}

const mutations = {
    setToken(state, token){
        state.token = token
    }
}

import maxiloVue from "maxilo-vue"
maxiloVue.make("store").add('kubernetes', {
    state,
    mutations,
    namespaced: true,
});