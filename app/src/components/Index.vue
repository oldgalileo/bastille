<template>
  <div class="hello">
    <!--<h1>Bastille</h1>-->
    <!--<h3>Memes</h3>-->
  </div>
</template>

<script>

export default {
  name: 'Index',
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

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1, h3 {
  font-weight: normal;
  margin: 0 auto;
}

h1 {
  font-size: 3rem;
  margin-bottom: 0;
}

ul {
  list-style-type: none;
  padding: 0;
}

li {
  display: inline-block;
  margin: 0 10px;
}

a {
  color: #42b983;
}
</style>
