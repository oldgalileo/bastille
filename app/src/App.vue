<template>
  <div id="app">
    <template v-if="isAuthed===true">
      <div class="navbar">
        <div class="content left">
          <div class="logo">Bastille</div>
        </div>
        <div class="content center">
          <div class="subtitle">{{ getQuote() }}</div>
        </div>
        <div class="content right">
          <div class="button">Welcome, {{ googleUser.getBasicProfile().getName() }}</div>
        </div>
      </div>
    </template>
    <template v-else>
      <div class="navbar">
        <div class="content left">
          <div class="logo">Bastille</div>
        </div>
        <div class="content center">
          <div class="subtitle">{{ getQuote() }}</div>
        </div>
        <div class="content right">
          <router-link to="/login" class="button" :callback="signIn">Sign In</router-link>
        </div>
      </div>
    </template>
    <router-view/>
  </div>
</template>

<script>
import gapi from 'gapi-client'

export default {
  name: 'app',
  data () {
    return {
      isAuthed: false,
      authError: '',
      googleAuthObject: null,
      googleUser: null,
      quotes: [
        'Patrick is already disappointed...',
        'In everything one thing is impossible: rationality.',
        'Inconceivable!',
        '"You only think I guessed wrong! That\'s what\'s so funny!" - Vizzini'
      ]
    }
  },
  mounted () {
    this.initialize()
  },
  methods: {
    getQuote: function () {
      return this.quotes[Math.floor(Math.random() * this.quotes.length)]
    },
    initialize: function () {
      gapi.load('client:oauth2', initClient.bind(this))
      function initClient () {
        gapi.client.init({
          'apiKey': 'Mm2sQnGVHJ1itd2-6lam9XWz',
          'clientId': '551253019639-slqo4k5kmueh8t09rrk6p4j5bltfnjf0.apps.googleusercontent.com',
          'scope': 'profile email'
        }).then(function () {
          this.googleAuthObject = gapi.auth2.getAuthInstance()
          if (this.googleAuthObject.currentUser.get().getBasicProfile()) {
            this.isAuthed = true
            this.googleUser = this.googleAuthObject.currentUser.get()
            this.initializeFormData()
          }
        })
      }
    },
    signIn: function () {
      this.authError = ''
      this.googleAuthObject.signIn().then(function () {
        this.updateSignInStatus(true)
      }.bind(this), function () {
        this.updateSignInStatus(false)
      }.bind(this))
    },
    updateSignInStatus: function (isAuthed) {
      if (isAuthed && this.googleAuthObject.currentUser.get().hasGrantedScopes('https://www.googleapis.com/auth/spreadsheets')) {
        if (this.googleAuthObject.currentUser.get().getBasicProfile().getEmail().split('@')[1] === 'nuevaschool.org') {
          this.googleUser = this.googleAuthObject.currentUser.get()
          this.isAuthed = true
        } else {
          this.googleAuthObject.disconnect()
          this.googleAuthObject.signOut()
          this.authError = 'You must be a member of the "nuevaschool.org" domain to use this application.'
          this.isAuthed = false
        }
      } else {
        this.googleAuthObject.signOut()
        this.isAuthed = false
      }
    }
  }
}
</script>

<style>
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 5vmin;
  padding: 5px 5px;
  display: flex;
  flex-flow: row wrap;
  justify-content: space-between;
  align-items: center;
  z-index: 1000;
  border-bottom: thin solid #BBB;
}

.navbar > .content {
  display: flex;
  flex-flow: row nowrap;
  align-items: center;
}

.navbar > .left > * {
  margin: 0 0 0 20px;
}

.navbar > .left > .divider {
  margin: 6px 0;
}

.navbar > .content > .logo {
  font-weight: normal;
  font-size: 2rem;
  padding-right: 20px;
}

.navbar > .right > * {
  margin: 0 20px 0 0;
}

.navbar > .content > .button {
  height: 4vmin;
  display: flex;
  justify-content: center;
  align-items: center;
  user-select: none;
  width: 80px;
  transition: all 0.2s ease;
  font-weight: bold;
  cursor: pointer;
  cursor: hand;
  border-bottom: 1px solid black;
}

.navbar > .content > .button:hover {
  border-bottom: 2px solid black;
}

#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  display: flex;
  justify-content: flex-start;
  width: 100vw;
  align-items: center;
  flex-flow: column nowrap;
}
</style>
