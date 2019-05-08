<template>
<span>{{humanReadable}}</span>
</template>

<script>
  let labels = ['B', 'KB', 'MB', 'GB', 'TB']
export default {
  props: [ 'bytes' ],
  computed: {
    humanReadable: function() {
      var index = 0;
      var obytes = this.bytes;
      while (obytes > 1024) {
        obytes /= 1024;
        index++;
      }
      if (index == 0) {
        return obytes + ' ' + labels[index];
      } else {
        return Math.round(obytes * 10) / 10 + ' ' + labels[index];
      }
    }
  }
}
</script>
