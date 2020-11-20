import maxiloVue from 'maxilo-vue'
import {
    Button,
    ConfigProvider,
    Card,
    Input,
    Checkbox,
    Statistic,
    InputNumber,
    Spin,
    Table,
    Modal,
    Select,
    Collapse,
    Tree,
    Icon,
    Divider,
    Message,
    Tabs,
    Notification,
    Popover,
    Tag,
    DatePicker,
    Descriptions,
    Switch,
    Result,
    Tooltip,
    Pagination,
    PageHeader,
    Breadcrumb,
    Badge,
    Radio,
    Skeleton,
    Drawer,
    Row,
    Col,
    Menu,
    Affix
} from 'ant-design-vue'
//Config-provider
maxiloVue.vue.use(Checkbox)
maxiloVue.vue.use(Statistic)
maxiloVue.vue.use(Button)
maxiloVue.vue.use(ConfigProvider)
maxiloVue.vue.use(Card)
maxiloVue.vue.use(Input)
maxiloVue.vue.use(InputNumber)
maxiloVue.vue.use(Spin)
maxiloVue.vue.use(Table)
maxiloVue.vue.use(Modal)
maxiloVue.vue.use(Select)
maxiloVue.vue.use(Tree)
maxiloVue.vue.use(Icon)
maxiloVue.vue.use(Divider)
maxiloVue.vue.use(Tabs)
maxiloVue.vue.use(Collapse)
maxiloVue.vue.use(Popover)
maxiloVue.vue.use(Tag)
maxiloVue.vue.use(DatePicker)
maxiloVue.vue.use(Descriptions)
maxiloVue.vue.use(Switch)
maxiloVue.vue.use(Result)
maxiloVue.vue.use(Tooltip)
maxiloVue.vue.use(Pagination)
maxiloVue.vue.use(PageHeader)
maxiloVue.vue.use(Breadcrumb)
maxiloVue.vue.use(Badge)
maxiloVue.vue.use(Radio)
maxiloVue.vue.use(Skeleton)
maxiloVue.vue.use(Drawer)
maxiloVue.vue.use(Row)
maxiloVue.vue.use(Col)
maxiloVue.vue.use(Menu)
maxiloVue.vue.use(Affix)
maxiloVue.vue.use(function(Vue){
    Vue.prototype.$message = Message
    Vue.prototype.$notification = Notification
    Vue.prototype.$info = Modal.info
    Vue.prototype.$success = Modal.success
    Vue.prototype.$error = Modal.error
    Vue.prototype.$warning = Modal.warning
    Vue.prototype.$warn = Modal.warn
    Vue.prototype.$confirm = Modal.confirm
})