import Vue from 'vue'
import App from './App.vue'
import Resource from 'vue-resource'

Vue.use(Resource)

Vue.config.productionTip = false

window.API = window.env_config.api + window.env_config.prefix + 'api';

Vue.http.options.root = window.API;

new Vue({
  render: h => h(App),
  mixins: [
    {
      beforeCreate: function() {
        if (!this.$options.http) {
          this.$options.http = this.$root.$options.http;
        }
      }
    }
  ]
}).$mount('#app')
