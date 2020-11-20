<script>
import cytoscape from "cytoscape"
import "codemirror/mode/dockerfile/dockerfile"
export default {
    render(){
        return <tool-curd httpKey="cloud" title="项目" 
            columns={this.cols}
            models={this.models}
            fetchUrl="/apis/deploy.hatobatsugu.gsc/v1/projects"
            fetchTransform={this.fetchTransform}
            vOn:add={this.onAdd}
            vOn:update={this.onAdd}
            vOn:vc={this.onVc}
            vOn:devAc={(item) => this.onAc(item, "dev")}
            vOn:prodAc={(item) => this.onAc(item, "prod")}
            vOn:delDev={this.onDelDev}></tool-curd>
    },
    data(){
        return {
            cols: [
                {field: "name", title: "项目"},
                {field: "git.url", title: "repo"},
                {field: "description", title: "描述"},
                {field: "resource.ports", type: "customer", title: "端口", hidden: true, option: {
                    customer_form: 'project-port'
                }, default: []},
                {field: "resource.configs", type: "customer", title: "配置", hidden: true, option: {
                    customer_form: 'project-config'
                }, default: []},
                {field: "resource.dockerfile", type: "code", title: "打包", hidden: true, option: {
                    language: 'dockerfile'
                }}
            ],
            models: [
                {key: "add", title: "新增", dispatchArea: "topBar",},
                {key: "update", title: "修改",},
            ]
        }
    },
    methods: {
        onAdd(data){
            let config = {
                apiVersion: "v1", kind: "ConfigMap",
                metadata: {name: ""}, data: {}
            }
            data.resource.configs.forEach(c => {
                let base = {
                    apiVersion: "v1", kind: "ConfigMap",
                    metadata: {
                        name: `project.${c.env}.${data.name}`,
                        labels: {
                            role: "app",
                            app: data.name,
                            env: c.env
                        }
                    }, data: {}
                }
                c.files.forEach(file => {
                    let paths = file.path.split("/")
                    base.data[paths[paths.length - 1]] = file.config
                })

                this.$kb.fullUpdateOrCreate(`/api/v1/namespaces/apps/configmaps`, base)
                    .catch(e => {
                        this.$notification.info({description: e.response.data.message})
                    })
            })

            this.$kb.fullUpdateOrCreate(`/apis/deploy.hatobatsugu.gsc/v1/namespaces/apps/projects`, {
                spec: data,
                apiVersion: "deploy.hatobatsugu.gsc/v1", kind: "Project",
                metadata: {
                    name: data.name,
                }
            })
                .catch(e => {
                    this.$notification.info({description: e.response.data.message})
                })
        },
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
            response.data.data = response.data.items.map(dp => dp.spec)
            response.data.total = response.data.data.length
            return response
        }
    }
}
</script>