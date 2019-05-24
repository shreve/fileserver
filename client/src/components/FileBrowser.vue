<template>
<div class="file-browser">
  <table>
    <thead>
      <tr>
        <th colspan="3" class="path">
          <path-breadcrumb :path="path" @change="handleBreadClick"></path-breadcrumb>
        </th>
      </tr>
      <tr>
        <th></th>
        <th>Name</th>
        <th>Size</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="file in files"
          :key="file.name"
          @click="handleClick($event, file)">
        <td><download-button :file="file" :path="stringPath"></download-button></td>
        <td>{{file.name}}<span v-if="file.dir">/</span></td>
        <td><file-size :bytes="file.size"></file-size></td>
      </tr>
    </tbody>
  </table>
</div>
</template>

<script>
import PathBreadcrumb from './PathBreadcrumb'
import FileSize from './FileSize'
import DownloadButton from './DownloadButton'

export default {
  components: {
    PathBreadcrumb,
    FileSize,
    DownloadButton
  },
  data: function() {
    return {
      path: [],
      files: [],
      cache: {},
      pending: false
    }
  },
  computed: {
    readPath: function() {
      let bits = decodeURIComponent(location.pathname).split('/').filter(String)
      let vbits = window.env_config.prefix.split('/').filter(String)
      while (bits[0] == vbits[0] && bits[0] != undefined) {
        bits.shift();
        vbits.shift();
      }
      return bits;
    },
    virtualPath: function() {
      return (window.env_config.prefix + this.path.join('/')).replace(/^\/+\//, '/');
    },
    apiPath: function() {
      return window.env_config.api + window.env_config.prefix + 'api';
    },
    stringPath: function() {
      return this.path.join('/')
    }
  },
  methods: {
    fetchList: function() {
      this.pending = true
      let params = { path: '/' + this.path.join('/') }
      this.$http
        .get('list', { params })
        .then(response => {
          this.files = response.body
          this.cache[this.path] = response.body
          this.pending = false
        })
    },
    handleClick: function(event, file) {
      if (file.dir) {
        if (this.pending) { return }
        this.path.push(file.name)
        this.visitPath();
      } else {
        let path = encodeURI(this.virtualPath + '/' + file.name).replace(/^\/+\//, '/')
        window.location = this.apiPath + '/download?path=' + path;
      }
    },
    handleBreadClick: function(path) {
      this.path = path;
      this.visitPath();
    },
    visitPath: function() {
      history.pushState({ path: this.path }, '', this.virtualPath)
      if (this.cache[this.path]) {
        this.files = this.cache[this.path];
        this.fetchList();
        this.pending = false;
      } else {
        this.fetchList();
      }
    },
    icon: function(file) {
      if (file.dir) { return 'D'; }
      return 'F';
    }
  },
  created: function() {
    this.path = this.readPath;

    this.fetchList();
    history.replaceState({ path: this.path }, location.pathname)

    window.onpopstate = (event) => {
      this.path = event.state.path;
      this.visitPath();
    }

    // setInterval(() => { this.fetchList() }, 5000);
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

td:first-child
  width: 0

</style>
