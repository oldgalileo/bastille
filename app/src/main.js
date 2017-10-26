// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import store from './store/index'
import auth from './auth'
import { sync } from 'vuex-router-sync'

Vue.config.productionTip = false
Vue.config.debug = true
Vue.config.devtools = true

sync(store, router)

/* eslint-disable no-new */
new Vue({
  el: '#app',
  store,
  router,
  template: '<App/>',
  components: { App },
  created: function () {
    auth.init()
  }
})
