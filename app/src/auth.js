import gapi from 'gapi-client'
import store from './store/index'
import { LOGIN_SUCCESS, LOGOUT } from './store/mutations'

export default {
  data: {
    user: null,
    instance: null
  },
  init: function () {
    gapi.load('client:auth2', loadClient.bind(this))
    function loadClient () {
      gapi.client.init({
        'apiKey': 'Mm2sQnGVHJ1itd2-6lam9XWz',
        'clientId': '551253019639-slqo4k5kmueh8t09rrk6p4j5bltfnjf0.apps.googleusercontent.com',
        'scope': 'profile email'
      }).then(function () {
        this.data.instance = gapi.auth2.getAuthInstance()
        if (this.data.instance.currentUser.get().getBasicProfile()) {
          this.data.user = this.data.instance.currentUser.get()
          store.commit(LOGIN_SUCCESS)
        } else {
          this.logout()
          store.commit(LOGOUT)
        }
      }.bind(this))
    }
  },
  login: function () {
    return this.data.instance.signIn().then(function () {
      return new Promise(function (resolve, reject) {
        // if (this.data.instance.currentUser.get().getBasicProfile().getEmail().split('@')[1] === 'nuevaschool.org') {
        this.data.user = this.data.instance.currentUser.get()
        resolve(this.data.user.getAuthResponse().id_token)
        // } else {
        //   reject('Must be a member of The Nueva School')
        // }
      }.bind(this))
    }.bind(this), function (error) {
      return Promise.reject(error)
    })
  },
  logout: function () {
    this.data.instance.disconnect()
    this.data.instance.signOut()
  }
}
