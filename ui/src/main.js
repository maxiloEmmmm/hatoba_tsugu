import maxiloVue from 'maxilo-vue'
import "./components/index"
import App from './App.vue'
import "./component"
import mvyu from 'maxilo-vue-ysz-ui/ysz-ui'
import 'antd-vue-tool/dist/index.css'
import tool from 'antd-vue-tool'
import "./validate"
import "./store"
import "./http"
import "./utils"
import "./route"

maxiloVue.vue.config.productionTip = true
maxiloVue.vue.use(mvyu)
maxiloVue.vue.use(tool)
maxiloVue.target = App
maxiloVue.run()
