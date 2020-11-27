<script>
const matchKeys = [
    "uri", "headers", "method", "queryParams"
]
const matchTypeKeys = [
    "prefix", "exact", "regex"
]
const istioStringMatch = {
    name: "istio-string-match",
    render(){
        return <tw-list-item1 fit index indexStart indexBorder items={[
            ...this.dataset.map((data, d_index) => {
                return {title: <a-tabs>
                    <a-button slot="tabBarExtraContent" size="small" vOn:click={() => this.onRemove(d_index)}>移除匹配组</a-button>
                    {Object.keys(data).map(d => {
                        return <a-tab-pane key={d} tab={d} force-render>
                            {data[d].map((item, d_index) => <ysz-list-item-top>
                                <a-button slot="top" size="small" vOn:click={() => this.onRemoveSearchTypeKey(data[d], d_index)}>移除匹配项</a-button>
                                <ysz-list row group={4}>
                                    <ysz-item-list>
                                        <span slot="left">key</span>
                                        <a-input size="small" value={item.key} vOn:change={e => this.onDataChange(e.target.value, item, "key")}/>
                                    </ysz-item-list>
                                    {matchTypeKeys.map(mtk => <ysz-item-list>
                                        <span slot="left">{mtk}</span>
                                        <a-input size="small" value={item[mtk]} vOn:change={e => this.onDataChange(e.target.value, item, mtk)}/>
                                    </ysz-item-list>)}
                                </ysz-list>
                            </ysz-list-item-top>)}
                            <a-button size="small" vOn:click={() => this.onNewSearchTypeKey(data[d])}>新增匹配项</a-button>
                        </a-tab-pane>}
                    )}
                </a-tabs>}
            }),
            {title: <a-button size={this.size} disabled={this.disabled} vOn:click={this.onNew}>新增匹配组</a-button>}
        ]}>
        </tw-list-item1>
    },
    props: {
        value: {type: Array, default(){
            return []
        }},
        disabled: {type: Boolean, default: false},
        size: {type: String, default: "small"}
    },
    methods: {
        onRemoveSearchTypeKey(arr, index) {
            arr.splice(index, 1)
            this.onChange()
        },
        onNewSearchTypeKey(match){
            let item = {key: ""}
            matchTypeKeys.forEach(mtk => {
                item[mtk] = ""
            })
            match.push(item)
            this.onChange()
        },
        onDataChange(v, d, k){
            d[k] = v
            this.onChange()
        },
        onChange(){
            this.$emit("change", this.dataset.map(r => {
                let matchs = {}

                Object.keys(r).forEach(rs => {
                    matchs[rs] = {}
                    r[rs].forEach(item => {
                        let index = matchTypeKeys.findIndex(key => !!item[key])
                        if(index == -1) {
                            return
                        }

                        let mk = matchTypeKeys[index]
                        matchs[rs][item.key] = {
                            [mk]: item[mk]
                        }
                    })
                    if(Object.keys(matchs[rs]).length == 0) {
                        delete matchs[rs]
                    }
                })
                return matchs
            }))
        },
        onNew(){
            let match = {}
            matchKeys.forEach(key => {
                match[key] = []
            })
            this.dataset.push(match)
        },
        onRemove(index) {
            this.dataset.splice(index, 1)
            this.onChange()
        }
    },
    data(){
        let rs = []
        // match组
        this.value.forEach(match => {
            let ms = {}
            // 以matchKey为key 组装match项 一个key有多个项
            matchKeys.forEach(key => {
                ms[key] = []
                // 如果存在key配置
                if(match[key]) {
                    // 循环项 组装item
                    Object.keys(match[key]).forEach(k => {
                        // 循环项中的匹配类型 组装item
                        Object.keys(match[key][k]).forEach(searchKey => {
                            let item = {
                                key: k,
                            }
                            // 将匹配类型展开 复制相应类型数据
                            matchTypeKeys.forEach(mtk => {
                                item[mtk] = searchKey == mtk ? match[key][k][searchKey] : ""
                            })
                            ms[key].push(item)
                        })
                        ms[key].push()
                    })
                }
            })
            rs.push(ms)
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
            ...this.dataset.map((data, d_index) => {
                return {title: <tw-list-item1 items={[
                        ...(!data.destination.host || !data.destination.port.number ? [{title: <a-alert message={`数据不全,不作为数据.`} type="warning" show-icon />}] : []),
                        {title: <ysz-list-item>
                            <span slot="left">目标</span>
                            <ysz-list-item-top>
                                <ysz-list-item slot="top">
                                    <span slot="left">✨ 服务</span>
                                    <tool-pick options={this.services} value={data.destination.host} vOn:change={(item) => this.onServiceChange(item, data)}><a-button size="small">{data.destination.host ? data.destination.host.replace(`-${this.env}.apps.svc.cluster.local`, "") : "选择"}</a-button></tool-pick>
                                </ysz-list-item>
                                <ysz-list-item>
                                    <span slot="left">✨ 端口</span>
                                    <tool-select options={this._ports} value={data.destination.port.number} vOn:change={v => data.destination.port.number = v}></tool-select>
                                </ysz-list-item>
                            </ysz-list-item-top>
                        </ysz-list-item>},
                        {title: <ysz-list-item>
                            <span slot="left">版本</span>
                            <tool-pick value={data.destination.subset} options={this.drs[data.destination.host]} vOn:change={(item) => this.onSubsetChange(item, data)}><a-button size="small">{data.destination.subset ? data.destination.subset : "选择"}</a-button></tool-pick>
                        </ysz-list-item>},
                        {title: <ysz-list-item>
                            <span slot="left">权重</span>
                            <a-input size="small" type="number" value={data.weight} vOn:change={this.onWeight}/>
                        </ysz-list-item>},
                        {title: <a-button size="small" vOn:click={() => this.onRemove(d_index)}>移除</a-button>}
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
        ports: {type: Array, default(){
            return []
        }},
        env: String,
        disabled: {type: Boolean, default: false},
        size: {type: String, default: "small"}
    },
    created(){this.fetch()},
    computed: {
        _ports(){
            return this.ports.map(port => {
                return {
                    label: port.name,
                    value: port.port
                }
            })
        }
    },
    methods: {
        onRemove(index){
            this.dataset.splice(index, 1)
            this.onChange()
        },
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
            this.$api.kubernetes.apis.deployment.list(`app=${app},env=${this.env}`).then(response => {
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
            this.$api.kubernetes.api.service.list(`env=${this.env}`)
                .then(response => {
                    this.services = response.data.items.map(service => {
                        let value = `${this.$utils.kbappid(service.metadata.labels.app, this.env)}.${this.$configs.cd_ns}.${this.$configs.cd_domain}`
                        this.fetchDR(service.metadata.labels.app, value)
                        return {
                            label: service.metadata.labels.app,
                            value
                        }
                    })
                })
        },
        onChange(){
            this.$emit("change", this.dataset.filter(d => !!d.destination.port.number && !!!!d.destination.host).map(d => {
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
            this.dataset.push({destination: {host: "", subset: "", port: {number: ""}}, weight: 0})
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
                    <span slot="left">Http</span>
                    <ysz-list-item-top>
                        <a-button slot="top" size="small" vOn:click={this.onNewHttp}>新增</a-button>
                        <a-collapse>
                            {this.https.map((http, index) => <a-collapse-panel key={http.name}>
                                <ysz-list-item slot="header">
                                    <ysz-list-item slot="left">
                                        <span slot="left">名称(唯一)</span><a-input value={http.name} size="small" vOn:click={evt => evt.stopPropagation()} vOn:change={evt => http.name = evt.target.value}/>
                                    </ysz-list-item>
                                    <a-space>
                                        <a-icon type="setting" vOn:click={(evt) => this.onRemoveHttp(evt, index)} />
                                        {index != 0 ? <a-icon type="arrow-up" vOn:click={(evt) => this.onSortHttp(evt, index, -1)}/> : null}
                                        {index != this.https.length - 1 ? <a-icon type="arrow-down" vOn:click={(evt) => this.onSortHttp(evt, index, 1)}/> : null}   
                                    </a-space>
                                </ysz-list-item>
                                <ysz-list>
                                    <ysz-list-item>
                                        <span slot="left">match</span>
                                        <istio-string-match value={http.match} vOn:change={v => http.match = v}></istio-string-match>
                                    </ysz-list-item>
                                    <ysz-list-item>
                                        <span slot="left">route</span>
                                        <istio-route value={http.route} env={this.env} vOn:change={v => http.route = v} ports={this.ports}></istio-route>
                                    </ysz-list-item>
                                    <ysz-list-item>
                                        <span slot="left">重写</span>
                                        <a-input size="small" value={http.rewrite.uri} vOn:change={e => http.rewrite.uri = e.target.value}></a-input>
                                    </ysz-list-item>
                                </ysz-list>
                            </a-collapse-panel>)}
                        </a-collapse>
                    </ysz-list-item-top>
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
            https: [],
            tls: [],
            loading: true,
            exist: false,
        }
    },
    created(){
        this.fetch()
    },
    props: {
        id: "",
        project: "",
        env: "",
        ports: {type: Array, default: () => []}
    },
    methods: {
        onSortHttp(evt, index, pos){
            evt.stopPropagation()
            let sort = this.https[index+pos]
            this.$set(this.https, index+pos, this.https[index])
            this.$set(this.https, index, sort)
        },
        onNewHttp(){
            this.https.push({
                name: this.$utils.tool.random("name-"),
                match: {},
                route: [],
                rewrite: {
                    uri: ""
                }
            })
        },
        onRemoveHttp(evt, index){
            evt.stopPropagation()
            this.https.splice(index, 1)
        },
        async end(){
            if(this.gateways.length == 0 || this.hosts.length == 0 || this.https.length == 0) {
                this.$message.info("资料不足")
                return 
            }

            this.https.forEach((http, indedx) => {
                if(http.route.length == 0) {
                    this.$message.info(`第${index+1}个http配置route不得为空`)
                    return 
                }
            })

            let item = {
                apiVersion: this.$api.kubernetes.apis.istio.vs.path.apiVersion(),
                kind: this.$api.kubernetes.apis.istio.vs.path.option.kind,
                metadata: {
                    name: this.$utils.kbappid(this.id, this.env),
                    namespace: "apps",
                    labels: {
                        app: this.id,
                        env: this.env,
                        role: this.$configs.role_app
                    }
                },
                spec: {
                    hosts: this.hosts,
                    gateways: this.gateways,
                    http: this.https.map(http => {
                        let tmp = {...http}
                        tmp.rewrite = {...tmp.rewrite}
                        if(!tmp.rewrite.uri) {
                            delete tmp.rewrite
                        }
                        return tmp
                    })
                }
            }

            this.$api.kubernetes.apis.istio.vs.fullUpdateOrCreate(item)
        },
        fetch(){
            Promise.all([
                this.$api.kubernetes.apis.istio.gateway.list()
                    .then(response => {
                        this.set.gateways = response.data.items.map(item => `${item.metadata.namespace}/${item.metadata.name}`)
                    }),
                ...(this.id ? [
                    this.$api.kubernetes.apis.istio.vs.get(this.$utils.kbappid(this.id, this.env))
                        .then(response => {
                            this.exist = true
                            this.gateways = response.data.spec.gateways
                            this.hosts = response.data.spec.hosts

                            this.https = response.data.spec.http === undefined ? [] : response.data.spec.http.map((h, index) => {
                                if(!h.name) {
                                    h.name = this.$utils.tool.random("name-")
                                }
                                if(!h.rewrite) {
                                    h.rewrite = {uri: ""}
                                }
                                return h
                            })
                            // this.tls = response.data.spec.tls === undefined ? [] : response.data.spec.tls
                        })
                        .catch(() => {})
                ] : [])
            ]).then(() => this.loading = false)
        }
    }
}
</script>