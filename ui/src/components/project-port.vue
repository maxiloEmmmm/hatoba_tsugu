<script>
// name	<string>
// port	<integer>
// protocol	<string>
// targetPort	<integer>
const DefaultPort = {
    port: 80,
    protocol: "TCP",
    name: "http",
    targetPort: 80
}
export default {
    name: "project-port",
    render(){
        return <tw-list-item1 index indexBorder items={[
            ...this.ports.map(port => {
                return {
                    title: <ysz-list-item-top>
                        <ysz-list-item slot="top">
                            <ysz-list-item slot="left">
                                <span slot="left">名字</span> <a-input size="small" value={port.name} vOn:change={e => this.onChange(port, 'name', e.target.value)}/>
                            </ysz-list-item>
                            <ysz-list-item>
                                <span slot="left">协议</span> <a-radio-group options={[
                                    {label: "TCP", value: "TCP"},
                                    {label: "UDP", value: "UDP"},
                                    {label: "SCTP", value: "SCTP"},
                                ]} size="small" value={port.protocol} vOn:change={e => this.onChange(port, 'protocol', e.target.value)}/>
                                
                            </ysz-list-item>
                        </ysz-list-item>
                        <ysz-list-item>
                            <ysz-list-item slot="left">
                                <span slot="left">服务端口</span> <a-input type="number" size="small" value={port.port} vOn:change={e => this.onChange(port, 'port', parseInt(e.target.value))}/>
                            </ysz-list-item>
                            <ysz-list-item>
                                <span slot="left">容器端口</span> <a-input type="number" size="small" value={port.targetPort} vOn:change={e => this.onChange(port, 'targetPort', parseInt(e.target.value))}/>
                            </ysz-list-item>
                        </ysz-list-item>
                    </ysz-list-item-top>
                }
            }),
            {title: <a-button size="small" vOn:click={this.onAdd}> 新增</a-button>}
        ]}></tw-list-item1>
    },
    props: {
        value: {
            type: Array, default: () => []
        }
    },
    watch: {
        value(val){
            this.ports = this.normalization(val)
        }
    },
    data(){
        return {
            ports: this.normalization(this.$props.value)
        }
    },
    methods: {
        onAdd(){
            this.ports.push({...DefaultPort})
            this.$emit("change", this.ports)
        },
        onChange(port, k, v){
            port[k] = v
            this.$emit("change", this.ports)
        },
        normalization(ports = []) {
            ports = Array.isArray(ports) ? ports : []
            return ports.map(port => {
                return {...port}
            })
        }
    }
}
</script>