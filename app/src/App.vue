<template>
  <div id="app">
      <div class="navbar">
        <div class="content left">
          <div class="logo">Bastille</div>
        </div>
        <!--<div class="content center">-->
          <!--<div class="subtitle">{{ getQuote() }}</div>-->
        <!--</div>-->
        <div class="content right">
          <div class="text">Welcome!</div>
          <div class="button">Upload</div>
          <div class="button">Leaderboard</div>
        </div>
      </div>
    <router-view/>
  </div>
</template>

<script>

export default {
  name: 'app',
  data () {
    return {
      quotes: [
        'Patrick is already disappointed...',
        'In everything one thing is impossible: rationality.',
        'Inconceivable!',
        '"You only think I guessed wrong! That\'s what\'s so funny!" - Vizzini'
      ]
    }
  },
  mounted () {},
  methods: {}
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  display: flex;
  justify-content: flex-start;
  width: 100vw;
  height: 100vh;
  align-items: center;
  flex-flow: column nowrap;
}

body {
  margin: 0;
}
</style>

<style scoped="true">
.navbar > * > .logo {
  font-weight: normal;
  font-size: 2rem;
  padding-right: 20px;
}

.navbar {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 5%;
  min-height: 30px;
  background: white;
  padding: 5px 5px;
  display: flex;
  flex-flow: row wrap;
  justify-content: space-between;
  align-items: center;
  z-index: 100;
  box-shadow: 0px 1px 0px 0px rgba(0, 0, 0, 0.5);
}

.navbar > .content {
  display: flex;
  flex-flow: row nowrap;
  align-items: center;
  height: 100%;
}

.navbar > .left > * {
  margin: 0 0 0 20px;
}

.navbar > .left > .divider {
  margin: 6px 0;
}

.navbar > .right > * {
  margin: 0 20px 0 0;
}

.navbar > .right {
  padding-right: 5px;
}


.navbar > .content > .text, .button {
  height: 80%;
  display: flex;
  justify-content: center;
  align-items: center;
  user-select: none;
  min-width: 80px;
  transition: all 0.2s ease;
  font-weight: bold;
  text-decoration: none;
  color: #2c3e50;
}

.navbar > .content > .button {
  font-weight: bold;
  cursor: pointer;
  cursor: hand;
  border-bottom: 1px solid black;
}

.navbar > .content > .button:hover {
  border-bottom: 2px solid black;
}
</style>

