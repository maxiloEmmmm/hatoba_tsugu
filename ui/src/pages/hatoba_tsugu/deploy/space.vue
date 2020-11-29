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
                    <tw-list-item2 items={pod.spec.containers.map(container => {
                        return {
                            title: <ysz-list-item>
                                <span slot="left">{container.name}</span>
                                <ysz-list row>
                                    <a-icon type="align-left" vOn:click={() => this.onContainerLog(pod.metadata.name, container.name)}/>
                                </ysz-list>
                            </ysz-list-item>,
                            desc: container.metric ? `${container.metric.usage.cpu} / ${container.metric.usage.memory}` : ""
                        }
                    })}></tw-list-item2>
                </a-collapse-panel>)}
            </a-collapse>
        </ysz-module-widget>
    },
    created(){
        this.renderPods()
        this.renderMetrics()
    },
    data(){
        return {
            set: {pods: [], podMetrics: []},
            runtime: {podSort: "cpu"}
        }
    },
    computed: {
        _pods() {
            return this.set.pods.map(pod => {
                let metricIndex = this.set.podMetrics.findIndex(metrics => metrics.metadata.name === pod.metadata.name)
                pod.spec.containers = pod.spec.containers.map(container => {
                    if(metricIndex === -1) {
                        container.metric = null
                    }else {
                        let containerMetricIndex = this.set.podMetrics[metricIndex].containers.findIndex(containerMetric => containerMetric.name === container.name)
                        container.metric = containerMetricIndex === -1 ? null : this.set.podMetrics[metricIndex].containers[containerMetricIndex]
                    }
                    return container
                })
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
    }
}
</script>