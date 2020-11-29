<script>
import CodeMirror from 'codemirror'

import 'codemirror/mode/shell/shell'
import 'codemirror/addon/hint/show-hint'
import 'codemirror/addon/hint/anyword-hint'
import 'codemirror/addon/edit/matchbrackets'

export default {
    name: "container-log",
    props: {
        pod: {type: String, default: ""},
        container: {type: String, default: ""},
        height: {type: String, default: "400px"},
    },
    render() {
        return <div ref="log" class="h-full"></div>
    },
    handler: null,
    ws: null,
    read: null,
    mounted(){
        this.connect()
    },
    destroyed(){this.close();},
    methods: {
        close(){
            if(this.handler) {
                this.clear()
            }
            this.ws && this.ws.close()
            this.ws = null
            this.handler = null
        },
        connect(){
            if(this.handler) {
                this.close()
                this.clear()
            }else {
                this.initWindow()
                this.initWs()
            }
        },
        initWindow(){
            this.handler = CodeMirror(this.$refs.log, {
                lineNumbers:true,
                theme:"seti",
                mode: 'shell',
                matchBrackets:true,
                smartIndent: false,
                lineWrapping: true,
                indentUnit: 4
            })
            this.handler.setSize("100%", this.height)
        },
        clear(){
            this.handler.setValue("")
        },
        initWs(){
            this.ws = new WebSocket(`ws://localhost:8000/cloud-ws/api/v1/namespaces/apps/pods/${this.pod}/log?container=${this.container}&stdin=true&stdout=true&tailLines=100&follow=true&tty=true`)
            this.ws.onopen = e => {
                console.log("open!", e)
                this.clear()
            };

            this.ws.onerror = e => {
                console.log("error", e)
            };
            this.ws.onclose = e => {
                this.close()
            };
            this.ws.onmessage = e => {
                let reader = new FileReader()
                reader.readAsText(e.data,"UTF-8")
                reader.onload = e => {
                    this.handler.setValue(`${this.handler.getValue()}${reader.result}`)
                    this.handler.setCursor(this.handler.lineCount(), 0)
                }
            };
        }
    }
}
</script>

<style lang="scss">
    @import '~codemirror/theme/seti.css';
    @import '~codemirror/lib/codemirror.css';
    @import '~codemirror/addon/hint/show-hint.css';
</style>