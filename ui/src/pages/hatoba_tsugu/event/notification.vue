<script>
import "codemirror/mode/markdown/markdown"
export default {
    render(){
        return <tool-curd httpKey="cloud" title="通知" 
            ref="curd"
            columns={this.cols}
            models={this.models}
            fetchUrl={this.$api.kubernetes.hatobatsugu.event.notification.url.multi}
            layout={[
                {key: "c2", col: 2},
                {key: "c1", col: 1},
            ]}
            fetchTransform={this.fetchTransform}
            vOn:add={this.onAdd}
            vOn:update={this.onAdd}/>
    },
    data(){
        return {
            cols: [
                {field: "name", title: "通知"},
                {field: "tpl", title: "模板", type: "code", option: {
                    language: "markdown"
                }, layout_key: "c1"},
                {field: "engine", title: "引擎", type: "select", option: {
                    selectOptions: () => this.set.channel
                }},
                {field: "filter", title: "过滤规则", hidden: true, type: "customer", option: {
                    customer_form: "notification-filter"
                }, layout_key: "c1"}
            ],
            models: [
                {key: "add", title: "新增", dispatchArea: "topBar",},
                {key: "update", title: "修改",},
            ],
            set: {
                channel: []
            }
        }
    },
    created(){
        this.fetchChannel()
    },
    methods: {
        onAdd(data){
            this.$api.kubernetes.hatobatsugu.event.notification.fullUpdateOrCreate({
                spec: data,
                apiVersion: this.$api.kubernetes.hatobatsugu.event.notification.path.apiVersion(), kind: this.$api.kubernetes.hatobatsugu.event.notification.path.option.kind,
                metadata: {
                    name: data.name,
                    namespace: this.$configs.cd_ns,
                }
            })
                .then(() => {
                    this.$refs.curd.refresh()
                    this.$api.event.refreshFilter()
                })
                .catch(e => {
                    this.$notification.info({description: e.response.data.message})
                })
        },
        fetchChannel(){
            this.$api.event.channel()
                .then(response => this.set.channel = response.data)
        },
        fetchTransform(response){
            response.data.data = response.data.items.map(dp => dp.spec)
            response.data.total = response.data.data.length
            return response
        }
    }
}
</script>