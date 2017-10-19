import Vuex from 'vuex'
import Vue from 'vue'
import auth from '../auth'
import { LOGIN, LOGIN_SUCCESS, LOGIN_FAILED, LOGOUT } from './mutations'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    authLoggedIn: false,
    authPending: false,
    authError: ''
  },
  mutations: {
    [LOGIN] (state) {
      state.authPending = true
      state.authError = ''
    },
    [LOGIN_SUCCESS] (state) {
      state.authPending = false
      state.authLoggedIn = true
    },
    [LOGIN_FAILED] (state, error) {
      state.authPending = false
      state.authError = error
    },
    [LOGOUT] (state) {
      state.authLoggedIn = false
    }
  },
  actions: {
    async login ({ commit }) {
      commit(LOGIN)
      await auth.login().then(function () {
        commit(LOGIN_SUCCESS)
      }, function (error) {
        commit(LOGIN_FAILED, error)
      })
    }
  }
})
