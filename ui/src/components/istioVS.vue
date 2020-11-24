<script>

const istioStringMatch = {
    name: "istio-string-match",
    render(){
        return <tw-list-item2 fit index indexStart indexBorder items={[
            ...this.dataset.map(data => {
                return {
                    title: <ysz-list-item-top>
                        <a-radio-group slot="top" size="small" value={data.type} vOn:change={e => this.onDataChange(e.target.value, data, "type")}>
                            <a-radio-button value="uri">Uri</a-radio-button>
                            <a-radio-button value="headers">Headers</a-radio-button>
                            <a-radio-button value="method">Method</a-radio-button>
                            <a-radio-button value="queryParams">QueryParams</a-radio-button>
                        </a-radio-group>
                        {data.type ? <a-button size={this.size} disabled={this.disabled} vOn:click={() => this.onNewKey(data)}>新增</a-button> : null}
                    </ysz-list-item-top>,
                    desc: data.items.map(d => {
                        return <ysz-list row group={4}>
                            <ysz-item-list>
                                <span slot="left">key</span>
                                <a-input size="small" value={d.key} vOn:change={e => this.onDataChange(e.target.value, d, "key")}/>
                            </ysz-item-list>
                            <ysz-item-list>
                                <span slot="left">exact</span>
                                <a-input size="small" value={d.exact} vOn:change={e => this.onDataChange(e.target.value, d, "exact")}/>
                            </ysz-item-list>
                            <ysz-item-list>
                                <span slot="left">prefix</span>
                                <a-input size="small" value={d.prefix} vOn:change={e => this.onDataChange(e.target.value, d, "prefix")}/>
                            </ysz-item-list>
                            <ysz-item-list>
                                <span slot="left">regex</span>
                                <a-input size="small" value={d.regex} vOn:change={e => this.onDataChange(e.target.value, d, "regex")}/>
                            </ysz-item-list>
                        </ysz-list>
                    })
                }
            }),
            {title: <a-button size={this.size} disabled={this.disabled} vOn:click={this.onNew}>新增</a-button>}
        ]}>
        </tw-list-item2>
    },
    props: {
        value: {type: Object, default(){
            return {}
        }},
        disabled: {type: Boolean, default: false},
        size: {type: String, default: "small"}
    },
    methods: {
        onDataChange(v, d, k){
            d[k] = v
            this.onChange()
        },
        onChange(){
            let obj = {}
            this.dataset.filter(r => !!r.type).forEach(r => {
                if(r.items.length ==  0) {
                    return
                }
                obj[r.type] =  {}
                r.items.forEach(item => {
                    let kinds = ["exact", "prefix", "regex"].filter(key => !!item[key])
                    if(kinds.length ==  0) {
                        return
                    }
                    kinds.forEach(t => {
                        obj[r.type][item.key][t] = item[t]
                    })
                })
            })
            this.$emit("change", obj)
        },
        onNew(){
            this.dataset.push({type: "", items: []})
        },
        onNewKey(d){
            d.items.push({key: "", exact: "", prefix: "", regex: ""})
        }
    },
    data(){
        let rs = []
        Object.keys(this.value).forEach(type => {
            let item = {type, items: []}
            Object.keys(this.value[type]).forEach(key => {
                item.items.push({key, ...this.value[type][key]})
            })
            rs.push(item)
        })
        return {
            newKey: '',
            dataset: Object.assign([], rs)
        }
    }
}

const istioRoute = {
    name: "istio-route",
    render(){
        return <tw-list-item1 fit index indexStart indexBorder items={[
            ...this.dataset.map(data => {
                return {
                    title: <tw-list-item1 items={[
                        {title: <tw-list-item1 items={[
                            {title: <ysz-list-item>
                                <span slot="left">目标</span>
                                <tool-pick options={this.services} value={data.destination.host} vOn:change={(item) => this.onServiceChange(item, data)}><a-button size="small">{data.destination.host ? data.destination.host.replace(`-${this.env}.apps.svc.cluster.local`, "") : "选择"}</a-button></tool-pick>
                            </ysz-list-item>},
                            {title: <ysz-list-item>
                                <span slot="left">版本</span>
                                <tool-pick value={data.destination.subset} options={this.drs[data.destination.host]} vOn:change={(item) => this.onSubsetChange(item, data)}><a-button size="small">{data.destination.subset ? data.destination.subset : "选择"}</a-button></tool-pick>
                            </ysz-list-item>},
                            {title: <ysz-list-item>
                                <span slot="left">权重</span>
                                <a-input size="small" type="number" value={data.weight} vOn:change={this.onWeight}/>
                            </ysz-list-item>}
                        ]}></tw-list-item1>}
                    ]}></tw-list-item1>
                }
            }),
            {title: <a-button size={this.size} disabled={this.disabled} vOn:click={this.onNew}>新增</a-button>}
        ]}>
        </tw-list-item1>
    },
    props: {
        value: {type: Array, default(){
            return []
        }},
        env: String,
        disabled: {type: Boolean, default: false},
        size: {type: String, default: "small"}
    },
    created(){this.fetch()},
    methods: {
        onWeight(e){
            data.weight = e.target.value
            this.onChange()
        },
        onSubsetChange(item, data){
            if(!item.value) {
                return
            }
            data.destination.subset = item.value
            this.onChange()
        },
        fetchDR(app, value){
            this.$kb.get(`/apis/apps/v1/namespaces/apps/deployments?labelSelector=${encodeURIComponent(`app=${app},env=${this.env}`)}`).then(response => {
                this.$set(this.drs, value, [])
                response.data.items.forEach(deployment => {
                    this.drs[value].push({
                        label: deployment.metadata.labels.version,
                        value: deployment.metadata.labels.version
                    })
                })
            })
        },
        onServiceChange(item, data){
            if(!item.value) {
                return
            }
            data.destination.host = item.value
            this.fetchDR(item.label)
            this.onChange()
        },
        fetch(){
            this.$kb.get(`/api/v1/namespaces/apps/services?labelSelector=${encodeURIComponent(`env=${this.env}`)}`)
                .then(response => {
                    this.services = response.data.items.map(service => {
                        let value = `app-${this.env}-${service.metadata.labels.app}.apps.svc.cluster.local`
                        this.fetchDR(service.metadata.labels.app, value)
                        return {
                            label: service.metadata.labels.app,
                            value
                        }
                    })
                })
        },
        onChange(){
            this.$emit("change", this.dataset.map(d => {
                let base = {
                    destination: {
                        ...d.destination
                    },
                    weight: d.weight
                }

                if (!base.weight) {
                    delete base.weight
                }
                if(!base.destination.subset) {
                    delete base.destination.subset
                }
                return base
            }))
        },
        onNew(){
            this.dataset.push({destination: {host: "", subset: ""}, weight: 0})
        }
    },
    data(){
        return {
            newKey: '',
            dataset: Object.assign([], this.value),
            drs: {},
            services: [],
        }
    }
}

export default {
    name: "istio-vs",
    components: {
        [istioStringMatch.name]: istioStringMatch,
        [istioRoute.name]: istioRoute
    },
    render(){
        return <a-skeleton loading={this.loading} active avatar>
            <ysz-list row group={1}>
                <ysz-list-item>
                    <span slot="left">域名</span>
                    <tool-tag value={this.hosts} vOn:change={v => this.hosts = v}></tool-tag>
                </ysz-list-item>
                <ysz-list-item>
                    <span slot="left">网关</span>
                    <a-checkbox-group value={this.gateways} options={this.set.gateways} vOn:change={v => this.gateways = v} />
                </ysz-list-item>
                <ysz-list-item>
                    <span slot="left">类型</span>
                    <a-radio-group value={this.type} options={["tls", "http"]} vOn:change={e => this.type = e.target.value} />
                </ysz-list-item>
                <ysz-list-item>
                    <span slot="left">{this.type}</span>
                    <ysz-list>
                        <ysz-list-item>
                            <span slot="left">match</span>
                            <istio-string-match value={this.http.match} vOn:change={v => this.http.match = v}></istio-string-match>
                        </ysz-list-item>
                        <ysz-list-item>
                            <span slot="left">route</span>
                            <istio-route value={this.http.route} env={this.env} vOn:change={v => this.http.route = v}></istio-route>
                        </ysz-list-item>
                    </ysz-list>
                </ysz-list-item>
                <ysz-list-item>
                    <a-button size="small" vOn:click={() => this.end()}>终了</a-button>
                </ysz-list-item>
            </ysz-list>
        </a-skeleton>
    },
    data(){
        return {
            set: {
                gateways: [],
            },
            hosts: [],
            gateways: [],
            http: {
                match: {},
                route: []
            },
            loading: true,
            type: "http",
            exist: false,
        }
    },
    created(){
        this.fetch()
    },
    props: {
        id: "",
        project: "",
        env: ""
    },
    methods: {
        async end(){
            if(this.gateways.length == 0 || this.hosts.length == 0 || this.http.route.length == 0) {
                this.$message.info("资料不足")
                return 
            }

            let item = {
                apiVersion: "networking.istio.io/v1beta1",
                kind: "VirtualService",
                metadata: {
                    name: `app-${this.env}-${this.id}`,
                    namespace: "apps"
                },
                spec: {
                    hosts: this.hosts,
                    gateways: this.gateways,
                    [this.type]: [this.http]
                }
            }

            if(!item.spec[this.type].match || Object.keys(item.spec[this.type].match) == 0) {
                item.spec[this.type] = [{route: item.spec[this.type][0].route}]
            }

            if(this.exist) {
                let response = await this.$kb.get(`/apis/networking.istio.io/v1beta1/namespaces/apps/virtualservices/${this.id}-${this.env}`)
                response.data.spec = item.spec
                this.$state.newState(this.$kb.put(`/apis/networking.istio.io/v1beta1/namespaces/apps/virtualservices/${response.data.metadata.name}`, response.data), {})
            }else {
                this.$state.newState(this.$kb.post(`/apis/networking.istio.io/v1beta1/namespaces/apps/virtualservices`, item), {})
            }
        },
        fetch(){
            Promise.all([
                this.$kb.get(`/apis/networking.istio.io/v1beta1/gateways`)
                    .then(response => {
                        this.set.gateways = response.data.items.map(item => `${item.metadata.namespace}/${item.metadata.name}`)
                    }),
                ...(this.id ? [
                    this.$kb.get(`/apis/networking.istio.io/v1beta1/namespaces/apps/virtualservices/app-${this.env}-${this.id}`)
                        .then(response => {
                            this.exist = true
                            this.gateways = response.data.spec.gateways
                            this.hosts = response.data.spec.hosts

                            this.type = response.data.spec.http === undefined ? "tls" : "http"
                            this.http = response.data.spec[this.type][0]
                        })
                        .catch(() => {})
                ] : [])
            ]).then(() => this.loading = false)
        }
    }
}
</script>