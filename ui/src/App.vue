<template>
  <div id="app">
    <router-view v-if="load"></router-view>
  </div>
</template>

<script>
export default {
  name: 'App',
  created(){
    this.$api.config.get().then(response => {
      this.$configs.merge(response.data)
      this.load = true
    })
  },
  data(){
    return {load: false}
  }
}
</script>

<style lang="scss">
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}
</style>
