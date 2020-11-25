<script>
import cytoscape from "cytoscape"
import "codemirror/mode/dockerfile/dockerfile"
export default {
    render(){
        return <tool-curd httpKey="cloud" title="项目" 
            columns={this.cols}
            models={this.models}
            fetchUrl={this.$api.kubernetes.hatobatsugu.deploy.project.url.multi}
            fetchTransform={this.fetchTransform}
            vOn:add={this.onAdd}
            vOn:update={this.onAdd}
            vOn:devAc={(item) => this.onAc(item, "dev")}
            vOn:prodAc={(item) => this.onAc(item, "prod")}
            layout={[
                {key: "c2", col: 2},
                {key: "c1", col: 1},
            ]}
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
                }, default: [], layout_key: "c1"},
                {field: "resource.configs", type: "customer", title: "配置", hidden: true, option: {
                    customer_form: 'project-config'
                }, default: [], layout_key: "c1"},
                {field: "resource.dockerfile", type: "code", title: "打包", hidden: true, option: {
                    language: 'dockerfile'
                }, layout_key: "c1"}
            ],
            models: [
                {key: "add", title: "新增", dispatchArea: "topBar",},
                {key: "update", title: "修改",},
                {key: "prodAc", title: "生产访问控制", type: "api", api: "event"},
                {key: "devAc", title: "开发访问控制", type: "api", api: "event"},
                {key: "delDev", title: "清除旧版本", type: "api", api: "event"}
            ]
        }
    },
    methods: {
        onAdd(data){
            let id = ""
            try {
                id = this.$utils.kbgitid(data.git.url)
            } catch (error) {
                alert(error)
                return
            }

            let config = {
                apiVersion: "v1", kind: "ConfigMap",
                metadata: {name: ""}, data: {}
            }
            data.resource.configs.forEach(c => {
                let base = {
                    apiVersion: "v1", kind: "ConfigMap",
                    metadata: {
                        name: this.$utils.kbappid(id, c.env),
                        labels: {
                            role: "app",
                            app: id,
                            env: c.env
                        }
                    }, data: {}
                }
                c.files.forEach(file => {
                    let paths = file.path.split("/")
                    base.data[paths[paths.length - 1]] = file.config
                })

                this.$api.kubernetes.api.configmap.fullUpdateOrCreate(base)
                    .catch(e => {
                        this.$notification.info({description: e.response.data.message})
                    })
            })

            this.$api.kubernetes.hatobatsugu.deploy.project.fullUpdateOrCreate({
                spec: data,
                apiVersion: "deploy.hatobatsugu.gsc/v1", kind: "Project",
                metadata: {
                    name: id,
                }
            })
                .catch(e => {
                    this.$notification.info({description: e.response.data.message})
                })
        },
        async onDelMore(item, table, env){
            let id = this.$utils.kbgitid(item.git.url)
            this.$api.kubernetes.apis.istio.vs.get(this.$utils.kbappid(id, env))
                .then(response => {
                    let useR = []
                    response.data.spec.http.forEach(h => {
                        h.route.forEach(r => {
                            useR.push(r.destination.subset)
                        })
                    })
                    if(useR.length > 0) {
                        this.$api.kubernetes.apis.istio.dr.get(this.$utils.kbappid(id, env))
                            .then(response => {
                                response.data.spec.subsets = response.data.spec.subsets.filter(s => useR.includes(s.name))
                                this.$state.newState(Promise.all([
                                    this.$api.kubernetes.apis.istio.dr.update(response.data),
                                    this.$api.kubernetes.apis.deployment.deleteLabel(`app=${id},env=${env},version notin (${useR.join(',')})`),
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
                    props: {id: this.$utils.kbgitid(item.git.url), env},
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