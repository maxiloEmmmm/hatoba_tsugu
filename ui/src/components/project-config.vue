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
    name: "project-config",
    render(){
        return <a-collapse>
            {this.env.map((e, index) => <a-collapse-panel key={e.env} header={e.env}>
                <a-space slot="extra">
                    <a-icon type="setting" vOn:click={(evt) => this.onAddConfig(evt, e, index)} />
                    <a-icon type="cloud-sync" vOn:click={(evt) => this.onSyncConfig(evt, e, index)} />
                </a-space>
                <a-collapse tabPosition="left">
                    {e.files.map((file, file_index) => <a-collapse-panel key={file_index}>
                        <ysz-list-item slot="header">
                            <ysz-list-item slot="left">
                                <span slot="left">全路径</span><a-input value={file.path} size="small" vOn:change={evt => this.onChange(file, "path", evt.target.value)}/>
                            </ysz-list-item>
                            <a-icon type="setting" vOn:click={(evt) => this.onRemoveConfig(evt, e.files, file_index)} />
                        </ysz-list-item>
                        <ysz-list no-line>
                            <ysz-list-item>
                                <span slot="left">描述</span><a-input value={file.description} size="small" vOn:change={evt => this.onChange(file, "description", evt.target.value)}/>
                            </ysz-list-item>
                            <ysz-list-item>
                                <span slot="left">配置</span><tool-code style="width:100%" value={file.config} vOn:change={evt => this.onChange(file, "config", evt)}></tool-code>
                            </ysz-list-item>
                        </ysz-list>
                    </a-collapse-panel>)}
                </a-collapse>
            </a-collapse-panel>)}
        </a-collapse>
    },
    props: {
        value: {
            type: Array, default: () => []
        }
    },
    watch: {
        value(val){
            this.env = this.normalization(val)
        }
    },
    data(){
        return {
            env: this.normalization(this.$props.value)
        }
    },
    methods: {
        onRemoveConfig(evt, files, index){
            evt.stopPropagation()
            files.splice(index, 1)
            this.$emit("change", this.env)
        },
        onAddConfig(evt, e){
            evt.stopPropagation()
            e.files.push({path: "", config: "", description: ""})
        },
        onSyncConfig(evt, e, index){
            evt.stopPropagation()
            this.env.forEach((env, e_index) => {
                if(e_index != index) {
                    let source_files = env.files.map(file => file.path)
                    env.files.push(...e.files.filter(file => !source_files.includes(file.path)))
                }
            })
        },
        onChange(file, k, v){
            file[k] = v
            this.$emit("change", this.env)
        },
        normalization(ports = []) {
            ports = Array.isArray(ports) ? ports : []
            if(ports.length == 0) {
                return this.normalization([
                    {env: "dev", files: []},
                    {env: "prod", files: []},
                ])
            }
            return ports.map(port => {
                return {...port}
            })
        }
    }
}
</script>