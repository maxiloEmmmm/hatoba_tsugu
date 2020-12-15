<script>
const DefaultPrometheus = {
    port: 80,
    path: "/metrics",
    enable: false
}
export default {
    name: "project-prometheus",
    render(){
        return <ysz-list-item>
            <a-checkbox vModel={this.enable} slot="left" vOn:change={this.onChange}>开启</a-checkbox>     
            <a-input-group compact>
                <a-input style="width: 30%" vModel={this.port} placeholder="端口" vOn:change={this.onChange}/>
                <a-input style="width: 70%" vModel={this.path} placeholder="路径" vOn:change={this.onChange}/>
            </a-input-group>
        </ysz-list-item>
    },
    props: {
        value: {
            type: Object, default: () => {return {...DefaultPrometheus}}
        }
    },
    data(){
        return this.normalization(this.$props.value)
    },
    methods: {
        onChange(){
            this.$emit("change", this.normalization(this.$data))
        },
        normalization(data) {
            data.port = data.port ? this.$utils.tool.number(data.port) : undefined
            return {...DefaultPrometheus, ...data}
        }
    }
}
</script>