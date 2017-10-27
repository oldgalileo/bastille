<template>
  <div class="navbar">
    <div class="content left">
      <div class="logo">Bastille</div>
    </div>
    <div class="content right">
      <div class="text">Welcome!</div>
      <template v-if="$store.state.authLoggedIn">
        <div class="button">Upload</div>
        <div class="button">Leaderboard</div>
        <div class="button">About</div>
        <div class="button" @click="signOut">Sign Out</div>
      </template>
      <template v-else>
        <div class="button">About</div>
        <div class="button" @click="signIn">Sign In</div>
      </template>
    </div>
  </div>
</template>

<script>
  import store from '../store/index'
  import router from '../router/index'

  export default {
    name: 'Navigation',
    methods: {
      signIn: () => {
        store.dispatch('login').then(() => {
          router.push('/dashboard')
        }, (reject) => {
          // Add a popup here to notify the user that authentication failed.
        })
      },
      signOut: () => {
        store.dispatch('logout').then(() => {
          router.push('/')
        })
      }
    }
  }
</script>

<style scoped="true">
  .navbar {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 50px;
    background-color: white;
    padding: 5px 5px;
    display: flex;
    flex-flow: row wrap;
    justify-content: space-between;
    align-items: center;
    box-shadow: 0 1px 3px rgba(0,0,0,0.12), 0 1px 2px rgba(0,0,0,0.24);
  }

  .content {
    display: flex;
    flex-flow: row nowrap;
    align-items: center;
    height: 100%;
  }

  .logo {
    margin-left: 20px;
    font-weight: bold;
    font-size: 2rem;
    padding-right: 20px;
  }

  .button, .text {
    margin-right: 20px;
  }

  .right {
    padding-right: 5px;
  }

  .button {
    height: 67%;
    display: flex;
    justify-content: center;
    align-items: center;
    user-select: none;
    transition: all 0.2s ease;
    text-decoration: none;
    color: #2c3e50;
    font-weight: bold;
    cursor: pointer;
    position: relative;
  }
  .button::after {
    height: 2px;
    width: 100%;
    position: absolute;
    display: block;
    content: ' ';
    bottom: 0;
    left: 0;
    background-color: #2c3e50;
    transition: transform 0.3s cubic-bezier(1, 0, 0, 1);;
  }
  .button:hover::after {
    transform: scaleY(1.5);
  }

</style>
