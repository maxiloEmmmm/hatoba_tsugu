<script>
export default {
    render(){
        return 1
    },
    created(){
        this.$kb.get("/api/v1/watch/namespaces/apps/pods", {params: {timeoutSeconds: 5, labelSelector: `app=lianbusinessplatformlianbigdatascreenplatform`}})
            .then(e => console.log(e))
    }
}
</script>