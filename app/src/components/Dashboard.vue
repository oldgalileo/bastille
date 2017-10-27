<template>
<div class="main">
  <section>
    <div class="content">
      <form method="POST" class="form" name="strategy" enctype="multipart/form-data">
        <ul class="inputs">
          <li class="input-line header-line"><label class="header">Submit Strategy</label></li>
          <li class="input-line">
            <label for="name">Name</label>
            <input type="text" id="name" name="name" placeholder="Tit-for-tat" required>
          </li>
          <li class="input-line">
            <label for="description">Description</label>
            <div class="textarea" id="description" placeholder="Tit-for-tat plays back the last move that your opponent made." contenteditable="true"></div>
            <!--<textarea id="description" name="description" maxlength="280" required rows="4"></textarea>-->
          </li>
          <li class="input-line">
            <label for="exec-button">File</label>
            <input type="file" class="file-input" id="exec" name="exec" hidden="true">
            <label for="exec" id="exec-button" class="file-label">Select File...</label>
          </li>
          <li class="input-line">
            <button type="submit">Submit</button>
          </li>
        </ul>
      </form>
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

  .form {
    min-width: 600px;
    width: 40%;
  }

  .inputs {
    width: 100%;
    display: flex;
    flex-flow: column nowrap;
    padding: 0;
  }

  .input-line {
    width: 100%;
    display: flex;
    flex-flow: row nowrap;
    justify-content: space-between;
    align-items: center;
    min-height: 75px;
    border-bottom: thin solid #2c3e50;
  }

  .header-line {
    justify-content: center;
  }

  .header-line .header {
    font-size: 2.0rem;
    font-weight: 500;
    padding-bottom: 50px;
  }

  .input-line label:not(.header) {
    font-size: 1.5rem;
    padding: 0 15px;
    border-right: 1px solid #2c3e50;
    line-height: 1.5rem;
    text-align: left;
    width: 25%;
  }

  .input-line input[type=text] {
    font-family: 'Avenir', Helvetica, Arial, sans-serif;
    height: 40px;
    border: none;
    -webkit-appearance: none;
    -moz-appearance: none;
    width: 100%;
    margin-left: 15px;
    font-size: 1.5rem;
  }

  .input-line input[type=text]:focus {
    outline: none;
    box-shadow: none;
    text-decoration-style: solid;
  }

  .input-line input[type=text]::placeholder {
    font-weight: 300;
    color: rgb(155, 155, 155);
  }

  .input-line .textarea {
    border: none;
    resize: none;
    width: 100%;
    margin: 15px 0px 15px 15px;
    font-size: 1.0rem;
    line-height: 1.15rem;
    min-height: 100px;
    text-align: left;
    font-weight: 400;
  }

  .input-line .textarea:focus {
    outline: none;
    box-shadow: none;
  }

  .input-line .textarea[placeholder]:empty:before {
    color: rgb(155, 155, 155);
    content: attr(placeholder);
    font-weight: 300;
  }

  .input-line .file-input {
    width: 0.1px;
    height: 0.1px;
    opacity: 0;
    overflow: hidden;
    position: absolute;
    z-index: -1;
    display: none;
  }

  .input-line .file-input:focus + label, .input-line .file-input + label:hover {

  }

  .input-line .file-label {
    height: 40px;
    width: 100%!important;
    border: none!important;
    display: flex;
    align-items: center;
    margin-left: 15px;
    padding: 0 !important;
    /*border-bottom: 1px solid #2c3e50!important;*/
    cursor: pointer;
  }

  .input-line:last-child {
    justify-content: center;
  }

  .input-line button[type=submit] {
    font-family: 'Avenir', Helvetica, Arial, sans-serif;
    color: #2c3e50;
    height: 40px;
    width: 40%;
    background: none;
    border: 1px solid #2c3e50;
    font-size: 1.5rem;
  }

  .input-line button[type=submit]:hover {
    background: #2c3e50;
    color: #ffffff;
  }

  .input-line button[type=submit]:active {
    background: #1e2b3c;
    outline: none;
  }

</style>
