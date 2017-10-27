<template>
  <div class="main">
    <div class="content">
      <h1>About</h1>
      <div class="writing">
        <p><strong>Bastille</strong> is an Iterated Prisoner's Dilemma (IPD) tournament manager designed to remove the limitations that afflict similar solutions. Bastille is language agnostic, which is achieved by <strong>containerizing the execution of uploaded strategies</strong>. Strategies interact with each other through an unshared supervisor which sends opponent move information via standard in and reads your program's moves via standard out.</p>
        <p>The simple nature of the interaction with the strategies allows for <strong>tens of thousands of rounds to happen in just a few seconds</strong>. Additionally, to ensure that the scale of the games isn't compromised by overly complicated strategies, there is a timeout on each move, after which your strategy will disqualified.</p>
        <p>This was made by Mickey, Leif, and Izzy. If you would like to become a contributor, feel free to email me (howard@getcoffee.io). If you want to check out the source, visit <a href="https://github.com/howardstark/bastille">https://github.com/howardstark/bastille</a>.</p>
      </div>
    </div>
  </div>
</template>

<script>
  export default {
    name: 'About'
  }
</script>

<style scoped="true">
  .writing {
    min-width: 600px;
    width: 40%;
    text-align: left;
    text-indent: 35px;
    line-height: 1.5rem;
  }

  a {
    color: #2c3e50;
  }
</style>
