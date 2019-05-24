<template>
<div class="button" :class="state" @click.stop="handleClick()" :title="title">
  {{text}}
</div>
</template>

<script>
const states = ["static", "pending", "ready"];
const texts = ["Compress", "Pending", "Download"];
const compressRate = 15 * 1024 * 1024;

export default {
  props: {
    file: {
    },
    path: String
  },
  data: function() {
    return {
      interval: null,
      stateNo: 0
    };
  },
  computed: {
    state: function() {
      return states[this.stateNo];
    },
    text: function() {
      return texts[this.stateNo];
    },
    cleanPath: function() {
      return encodeURI('/' + this.path + '/' + this.file.name).replace(/^\/+\//, '/');
    },
    compressPeriod: function() {
      return Math.round(this.file.size / compressRate)
    },
    title: function() {
      if (this.stateNo === 2) { return }
      let est = this.compressPeriod;
      if (est < 1) { return "Basically instantly"; }
      return "About " + est + "s";
    }
  },
  methods: {
    handleClick: function() {
      switch(this.state) {
      case "static":
        if (this.confirmCompress()) {
          this.startCompression();
        }
        break;
      case "pending":
        break;
      case "ready":
        this.download();
      }
    },

    confirmCompress: function() {
      let est = this.compressPeriod;
      if (est >= 60) {
        return confirm("Are you sure you want to compress this folder? It will take about " + est + "s");
      }
      return true;
    },

    startCompression: function() {
      let params = { path: this.cleanPath };
      this.$http
        .get('download', { params })
        .then(() => {}, () => {
          this.startWatcher();
          this.stateNo = 1;
        })
    },

    startWatcher: function() {
      this.interval = setInterval(() => {
        this.checkStatus();
      }, 500)
    },

    checkStatus: function() {
      let params = { path: this.cleanPath };
      this.$http
        .get('status', { params })
        .then(() => {
          this.stateNo = 2;
          clearInterval(this.interval);
        })
    },

    download: function() {
      window.location = window.API + '/download?path=' + this.cleanPath;
    }
  },
  created: function() {
    if (this.file.dir) {
      this.stateNo = (this.file.zipped ? 2 : 0);
    } else {
      this.stateNo = 2;
    }
  }
}
</script>

<style lang="sass">
.button
  border: none
  border-radius: 4px
  padding: 0.2em 0.5em 0.3em
  color: white
  display: inline-block

.static
  background: darkgrey

.pending
  background: gold

.ready
  background: dodgerblue
</style>
