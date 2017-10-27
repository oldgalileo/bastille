<template>
  <div id="app">
    <Navigation></Navigation>
    <router-view/>
  </div>
</template>

<script>
  import Navigation from './components/Navigation'

  export default {
    name: 'app',
    components: {Navigation},
    data () {
      return {
        quotes: [
          'Patrick is already disappointed...',
          'In everything one thing is impossible: rationality.',
          'Inconceivable!',
          '"You only think I guessed wrong! That\'s what\'s so funny!" - Vizzini'
        ]
      }
    }
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

.main {
  height: 100%;
  width: 100%;
  display: flex;
  flex-flow: row wrap;
  justify-content: space-between;
  align-items: center;
}

section {
  width: 100%;
  min-height: 100%;
  display: flex;
  flex-flow: row wrap;
  justify-content: center;
}

.content {
  width: 100%;
  min-height: 100%;
  display: flex;
  flex-flow: column nowrap;
  align-items: center;
  justify-content: center;
}
</style>

