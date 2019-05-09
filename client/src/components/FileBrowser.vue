<template>
<div class="file-browser">
  <table>
    <thead>
      <tr>
        <th colspan="2" class="path">
          <path-breadcrumb :path="path" @change="handleBreadClick"></path-breadcrumb>
        </th>
      </tr>
      <tr>
        <th>Name</th>
        <th>Size</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="file in files"
          :key="file.name"
          @click="handleClick(file)">
        <td>{{file.name}}<span v-if="file.dir">/</span></td>
        <td><span v-if="!file.dir"><file-size :bytes="file.size"></file-size></span></td>
      </tr>
    </tbody>
  </table>
</div>
</template>

<script>
  import PathBreadcrumb from './PathBreadcrumb'
  import FileSize from './FileSize'

export default {
  components: {
    PathBreadcrumb,
    FileSize
  },
  data: function() {
    return {
      path: [],
      files: []
    }
  },
  computed: {
    readPath: function() {
      let bits = location.pathname.split('/').filter(String)
      let vbits = window.env_config.prefix.split('/').filter(String)
      while (bits[0] == vbits[0] && bits[0] != undefined) {
        bits = bits.shift();
        bits = vbits.shift();
      }
      return bits;
    },
    virtualPath: function() {
      return window.env_config.prefix + this.path.join('/');
    },
    apiPath: function() {
      return window.env_config.prefix + 'api';
    }
  },
  methods: {
    fetchList: function() {
      let params = { path: '/' + this.path.join('/') }
      this.$http
        .get(this.apiPath + '/list', { params })
        .then(response => {
          this.files = response.body
        })
    },
    handleClick: function(file) {
      if (file.dir) {
        this.path.push(file.name)
        this.visitPath();
      } else {
        window.location = window.env_config.prefix + '__files/' + this.path.join('/') + '/' + file.name;
      }
    },
    handleBreadClick: function(path) {
      this.path = path;
      this.visitPath();
    },
    visitPath: function() {
      console.log(this.virtualPath);
      history.pushState({ path: this.path }, '', this.virtualPath)
      this.fetchList();
    }
  },
  created: function() {
    this.path = this.readPath;

    this.fetchList();
    history.replaceState({ path: this.path }, location.pathname)

    window.onpopstate = (event) => {
      this.path = event.state.path;
      this.fetchList();
    }

    setInterval(() => { this.fetchList() }, 5000);
  }
}
</script>

<style lang="sass">
  table
    text-align: left
    width: 100%
    max-width: 50em
    border: 1px solid #ddd
    margin: 0 auto
    border-collapse: collapse

td, th
  padding: 0.5em
  border: 0
  margin: 0

tbody tr:hover
  background: #f5f5f5
  cursor: pointer

.path
  border-bottom: 1px solid #ddd

</style>
