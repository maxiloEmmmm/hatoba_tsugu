<script>
export default {
    render(){
        return <ysz-module-widget>
            <ysz-list-item slot="title">
                <span slot="left">工作节点</span>
                <ysz-list row>
                    <ysz-list-item>
                        <span slot="left">排序</span>
                        <a-radio-group vModel={this.runtime.podSort} size="small" options={[
                            {label: "cpu", value: "cpu"},
                            {label: "mem", value: "mem"}
                        ]}/>
                    </ysz-list-item>
                </ysz-list>
            </ysz-list-item>
            <a-collapse>
                {this._pods.map(pod => <a-collapse-panel key={pod.metadata.name} header={pod.metadata.name}>
                    <span slot="extra">{`${pod.metric.view.cpu}/${pod.metric.view.cpuUsePercent}/${pod.metric.view.cpuPercent}, ${pod.metric.view.mem}/${pod.metric.view.memUsePercent}/${pod.metric.view.memPercent}`} | {pod.spec.nodeName}</span>
                    <ysz-list>
                        <ysz-module-card title="容器">
                            <ysz-list row groups={4}>
                                {pod.spec.containers.map(container => <ysz-module-widget>
                                    <ysz-list-item-top slot="title">
                                        <span slot="top">{container.name}</span>
                                        <ysz-list row>
                                            <a-icon type="align-left" vOn:click={() => this.onContainerLog(pod.metadata.name, container.name)}/>
                                        </ysz-list>
                                    </ysz-list-item-top>
                                    <ysz-module-widget title="资源占用">
                                        <tw-list-item2 items={[{
                                            title: "用量(cpu/内存)",
                                            desc: `${container.metric.view.cpu} / ${container.metric.view.mem}`
                                        }, {
                                            title: "已使用占比(宿主机cpu/内存)",
                                            desc: `${container.metric.view.cpuUsePercent} / ${container.metric.view.memUsePercent}`
                                        }, {
                                            title: "占比(宿主机cpu/内存)",
                                            desc: `${container.metric.view.cpuPercent} / ${container.metric.view.memPercent}`
                                        }]}></tw-list-item2>
                                    </ysz-module-widget>
                                </ysz-module-widget>)}
                            </ysz-list>
                        </ysz-module-card>
                        <ysz-module-card title={`所在主机${pod.spec.nodeName}`}>
                            <tw-list-item2 items={[{
                                title: "用量(cpu/内存)",
                                desc: `${pod.nodeView.use.cpu} / ${pod.nodeView.use.mem}`
                            }]}></tw-list-item2>
                        </ysz-module-card>
                    </ysz-list>
                </a-collapse-panel>)}
            </a-collapse>
        </ysz-module-widget>
    },
    created(){
        this.renderPods()
        this.renderMetrics()
        this.renderNodes()
        this.fetchNode()
    },
    data(){
        return {
            set: {pods: [], nodes: [], podMetrics: [], nodeMetrics: []},
            runtime: {podSort: "cpu"}
        }
    },
    computed: {
        _pods() {
            return this.set.pods.map(pod => {
                let nodeIndex = this.set.nodeMetrics.findIndex(node => node.metadata.name === pod.spec.nodeName)
                pod.nodeMetric = nodeIndex === -1 ? null : this.set.nodeMetrics[nodeIndex]
                
                let node = this.set.nodes.filter(node => node.metadata.name == pod.spec.nodeName)[0]
                pod.nodeView = {
                    use: {
                        cpu: pod.nodeMetric && node ? `${this.$utils.percent((parseInt(pod.nodeMetric.usage.cpu) / (1000 * 1000)) / (parseInt(node.status.capacity.cpu) * 1000))}%, ${node.status.capacity.cpu}C` : "",
                        mem: pod.nodeMetric && node ? `${this.$utils.percent(parseInt(pod.nodeMetric.usage.memory) / parseInt(node.status.capacity.memory))}%, ${this.$utils.ktoview(parseInt(node.status.capacity.memory))}` : "",
                    }
                }

                let metricIndex = this.set.podMetrics.findIndex(metrics => metrics.metadata.name === pod.metadata.name)
                pod.metric = {
                    cpu: 0,
                    mem: 0,
                }
                pod.spec.containers = pod.spec.containers.map(container => {
                    if(metricIndex === -1) {
                        container.metric = null
                    }else {
                        let containerMetricIndex = this.set.podMetrics[metricIndex].containers.findIndex(containerMetric => containerMetric.name === container.name)
                        container.metric = containerMetricIndex === -1 ? null : this.set.podMetrics[metricIndex].containers[containerMetricIndex]

                        if(container.metric) {
                            pod.metric.cpu += parseInt(container.metric.usage.cpu)
                            pod.metric.mem += parseInt(container.metric.usage.memory)
                        }
                        let pm = parseInt(parseInt(container.metric.usage.cpu) / (1000 * 1000))
                        container.metric.view = {
                            cpu: `${pm}m`,
                            mem: this.$utils.ktoview(parseInt(container.metric.usage.memory)),
                            cpuUsePercent: pod.nodeMetric && node ? `${this.$utils.percent(parseInt(container.metric.usage.cpu) / parseInt(pod.nodeMetric.usage.cpu))}%` : "",
                            memUsePercent: pod.nodeMetric && node ? `${this.$utils.percent(parseInt(container.metric.usage.memory) / parseInt(pod.nodeMetric.usage.memory))}%` : "",
                            cpuPercent: pod.nodeMetric && node ? `${this.$utils.percent(pm / (parseInt(node.status.capacity.cpu) * 1000))}%` : "",
                            memPercent: pod.nodeMetric && node ? `${this.$utils.percent(parseInt(container.metric.usage.memory) / parseInt(node.status.capacity.memory))}%` : "",
                        }
                    }
                    return container
                })

                let ppm = parseInt(parseInt(pod.metric.cpu) / (1000 * 1000))
                pod.metric.view = {
                    cpu: `${ppm}m`,
                    mem: this.$utils.ktoview(parseInt(pod.metric.mem)),
                    cpuUsePercent: node ? `${this.$utils.percent(parseInt(pod.metric.cpu) / parseInt(pod.nodeMetric.usage.cpu))}%` : "",
                    memUsePercent: node ? `${this.$utils.percent(parseInt(pod.metric.mem) / parseInt(pod.nodeMetric.usage.memory))}%` : "",
                    cpuPercent: node ? `${this.$utils.percent(ppm / (parseInt(node.status.capacity.cpu) * 1000))}%` : "",
                    memPercent: node ? `${this.$utils.percent(pod.metric.mem / parseInt(node.status.capacity.memory))}%` : "",
                }
                return pod
            })
        }
    },
    methods: {
        onContainerLog(pod, container){
            const h = this.$createElement;
            this.$info({
                title: '日志',
                width: "80%",
                content: h('container-log', {props: {
                    pod,
                    container,
                    height: "400px"
                }}),
                onOk() {},
            });
        },
        renderPods() {
            this.$utils.setInterval(() => this.$api.kubernetes.api.pod.list(`app=${this.$route.params.id}`).
                then(response => this.set.pods = response.data.items), 5)
        },
        renderMetrics() {
            this.$utils.setInterval(() => this.$api.kubernetes.apis.metrics.pod.list(`app=${this.$route.params.id}`).
                then(response => this.set.podMetrics = response.data.items), 5)
        },
        renderNodes() {
            this.$utils.setInterval(() => this.$api.kubernetes.apis.metrics.node.list().
                then(response => this.set.nodeMetrics = response.data.items), 5)
        },
        fetchNode() {
            this.$api.kubernetes.api.node.list().
                then(response => this.set.nodes = response.data.items)
        }
    }
}
</script>