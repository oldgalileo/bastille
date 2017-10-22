<template>
<div class="main">
  <section>
    <div class="content">
      <h1>Dashboard</h1>
    </div>
  </section>
</div>
</template>

<script>
  export default {
    name: 'Dashboard'
  }
</script>

<style scoped>
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
