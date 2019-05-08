<template>
<nav class="breadcrumb">
<ul>
  <li @click="handleClick(0)">root</li>
  <li v-for="(crumb, i) in path"
      :key="i"
      @click="handleClick(i+1)">{{crumb}}</li>
</ul>
</nav>
</template>

<script>
export default {
  props: [
    'path'
  ],
  methods: {
    handleClick: function(i) {
      this.$emit('change', this.path.slice(0, i));
    }
  }
}
</script>

<style lang="sass">
.breadcrumb
  font-weight: normal
  font-family: monospace

  ul
    list-style-type: none
    margin: 0
    padding: 0

li
  display: inline-block
  cursor: pointer

  &:last-child
    cursor: inherit

  &:after
    content: '/'
    display: inline-block
</style>
