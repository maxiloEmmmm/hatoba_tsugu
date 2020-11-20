<script>
import cytoscape from "cytoscape"
export default {
    render(){
        return <tool-curd httpKey="cloud" title="项目" 
            columns={this.cols}
            models={this.models}
            fetchUrl={`/api/v1/namespaces/apps/services?labelSelector=${encodeURIComponent("role=app")}`}
            fetchTransform={this.fetchTransform}
            vOn:vc={this.onVc}
            vOn:devAc={(item) => this.onAc(item, "dev")}
            vOn:prodAc={(item) => this.onAc(item, "prod")}
            vOn:delDev={this.onDelDev}></tool-curd>
    },
    data(){
        return {
            cols: [
                {field: "name", title: "项目"},
            ],
            models: [
                {key: "prodAc", title: "生产访问控制", type: "api", api: "event"},
                {key: "devAc", title: "开发访问控制", type: "api", api: "event"},
                {key: "delDev", title: "清除旧版本", type: "api", api: "event"}
            ]
        }
    },
    methods: {
        async onDelMore(item, table, env){
            this.$kb.get(`/apis/networking.istio.io/v1beta1/namespaces/apps/virtualservices/${item.name}-${env}`)
                .then(response => {
                    let useR = []
                    response.data.spec.http.forEach(h => {
                        h.route.forEach(r => {
                            useR.push(r.destination.subset)
                        })
                    })
                    if(useR.length > 0) {
                        this.$kb.get(`/apis/networking.istio.io/v1beta1/namespaces/apps/destinationrules/${item.name}-${env}`)
                            .then(response => {
                                response.data.spec.subsets = response.data.spec.subsets.filter(s => useR.includes(s.name))
                                this.$state.newState(Promise.all([
                                    this.$kb.put(`/apis/networking.istio.io/v1beta1/namespaces/apps/destinationrules/${item.name}-${env}`, response.data),
                                    this.$kb.delete(`/apis/apps/v1/namespaces/apps/deployments?labelSelector=${encodeURIComponent(`app=${item.name},env=${env},version notin (${useR.join(',')})`)}`),
                                ]), {})
                            })
                    }
                })
        },
        async onDelDev({item, table}) {
            this.onDelMore(item, table, "dev")
            // this.onDelMore(item, table, "prod")
        },
        async onAc({item}, env){
            const h = this.$createElement
            let modal = this.$info({
                title: `${env}-访问流`,
                width: 1024,
                footer: null,
                content: h('istio-vs', {
                    props: {id: item.name, env},
                    on: {done: () => {
                        modal.destroy()
                        table.refresh()
                    }}
                })
            })
        },
        async onVc({item, table}){
            const h = this.$createElement
            let modal = this.$info({
                title: '版本',
                width: 1024,
                content: h('business-app-traffic', {
                    props: {data: item},
                    on: {done: () => {
                        modal.destroy()
                        table.refresh()
                    }}
                })
            })
        },
        fetchTransform(response){
            let project = {}
            response.data.items.forEach(p => {
                if(!project[p.metadata.labels.app]) {
                    project[p.metadata.labels.app] = {
                        services: [p],
                        name: p.metadata.labels.app,
                    }
                }else {
                    project[p.metadata.labels.app].services.push(p)
                }
            })
            response.data.data = Object.values(project)
            response.data.total = response.data.data.length
            return response
        }
    }
}
</script>