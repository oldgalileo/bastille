import Vuex from 'vuex'
import Vue from 'vue'
import auth from '../auth'
import { LOGIN, LOGIN_SUCCESS, LOGIN_FAILED, LOGOUT, LOADED } from './mutations'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    loading: true,
    authLoggedIn: false,
    authPending: false,
    authError: '',
    authToken: ''
  },
  mutations: {
    [LOADED] (state) {
      state.loading = false
    },
    [LOGIN] (state) {
      state.authPending = true
      state.authError = ''
    },
    [LOGIN_SUCCESS] (state, token) {
      state.authPending = false
      state.authLoggedIn = true
      state.authToken = token
    },
    [LOGIN_FAILED] (state, error) {
      state.authPending = false
      state.authError = error
    },
    [LOGOUT] (state) {
      state.authLoggedIn = false
      state.authToken = ''
    }
  },
  actions: {
    async login ({ commit }) {
      commit(LOGIN)
      return auth.login().then(function (token) {
        commit(LOGIN_SUCCESS, token)
      }, function (error) {
        commit(LOGIN_FAILED, error)
      })
    },
    logout ({ commit }) {
      auth.logout()
      commit(LOGOUT)
    }
  }
})
