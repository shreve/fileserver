<template>
  <div id="app">
    <h1>Files</h1>
    <file-browser></file-browser>
  </div>
</template>

<script>
  import FileBrowser from './components/FileBrowser.vue'

  export default {
  components: {FileBrowser}
  }
</script>

<style lang="scss">
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
</style>
