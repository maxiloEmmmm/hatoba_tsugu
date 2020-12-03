<script>
const DefaultFilter = {
    type: "",
    key: "",
    val: ""
}
export default {
    name: "notification-filter",
    render(){
        return <tw-list-item1 index indexBorder items={[
            ...this.filters.map(filter => {
                return {
                    title: <a-input-group compact size="small">
                        <a-select size="small" options={this.set.filters} vModel={filter.type} vOn:change={this.onChange}></a-select>
                        <a-input style="width: 20%" placeholder="key" vModel={filter.key} vOn:change={this.onChange}/>
                        <a-input style="width: 50%" placeholder="val" vModel={filter.val} vOn:change={this.onChange}/>
                    </a-input-group>
                }
            }),
            {title: <a-button size="small" vOn:click={this.onAdd}>{!this.runtime.loading ? "新增" : "加载中..."}</a-button>}
        ]}></tw-list-item1>
    },
    props: {
        value: {
            type: Array, default: () => []
        }
    },
    watch: {
        value(val){
            this.filters = this.normalization(val)
        }
    },
    created(){
        this.fetchType()
    },
    data(){
        return {
            filters: this.normalization(this.$props.value),
            set: {
                filters: []
            },
            runtime: {
                loading: true
            }
        }
    },
    methods: {
        fetchType(){
            this.$api.event.filter()
                .then(response => {
                    this.set.filters = response.data
                    this.runtime.loading = false
                })
        },
        onAdd(){
            let filter = {...DefaultFilter}
            if(this.set.filters.length > 0) {
                filter.type = this.set.filters[0].value
            }
            this.filters.push(filter)
            this.$emit("change", this.filters)
        },
        onChange(){
            this.$emit("change", this.filters)
        },
        normalization(filters = []) {
            filters = Array.isArray(filters) ? filters : []
            return filters.map(filter => {
                return {...filter}
            })
        }
    }
}
</script>