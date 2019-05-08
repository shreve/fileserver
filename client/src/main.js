import Vue from 'vue'
import App from './App.vue'
import Resource from 'vue-resource'

Vue.use(Resource)

Vue.config.productionTip = false

new Vue({
  render: h => h(App)
}).$mount('#app')
