<template>
<div class="main">
  <div class="overlay" v-if="!isInitial">
    <div class="dialog">
      <div class="content" v-if="isUploading">
        <div class="loader">
          <div class="inside">Uploading</div>
          <div class="circle"></div>
        </div>
      </div>
      <div class="content" v-else-if="isSuccess">
        <h2>Success!</h2>
        <p>Your strategy was uploaded successful.</p>
        <button @click="reset()">Close</button>
      </div>
      <div class="content" v-else-if="isFailed">
        <h2>Uh-oh!</h2>
        <p>Could not upload strategy...</p>
        <button @click="reset()">Close</button>
      </div>
    </div>
  </div>
  <section>
    <div class="content">
      <form method="POST" class="form" name="strategy" enctype="multipart/form-data" @submit.prevent="uploadStrategy()">
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
            <input type="file" class="file-input" id="exec" name="exec" hidden="true" @change="handleFile( $event.target.files )">
            <template v-if="strategyFileName === ''">
              <label for="exec" id="exec-button" class="file-label">Select File...</label>
            </template>
            <template v-else>
              <label for="exec" id="exec-button" class="file-label">{{ strategyFileName }}</label>
            </template>
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
  import * as config from '@/config'
  import auth from '@/auth'

  const STATUS_INITIAL = 0
  const STATUS_UPLOADING = 1
  const STATUS_SUCCESS = 2
  const STATUS_FAILED = 3

  export default {
    name: 'Dashboard',
    data () {
      return {
        status: STATUS_INITIAL,
        strategyFile: null,
        strategyFileName: '',
        uploadResponse: ''
      }
    },
    computed: {
      isInitial () {
        return this.status === STATUS_INITIAL
      },
      isUploading () {
        return this.status === STATUS_UPLOADING
      },
      isSuccess () {
        return this.status === STATUS_SUCCESS
      },
      isFailed () {
        return this.status === STATUS_FAILED
      }
    },
    methods: {
      reset () {
        this.status = STATUS_INITIAL
        this.strategyFile = null
        this.strategyFiileName = ''
        this.uploadResponse = ''
      },
      handleFile (files) {
        this.strategyFileName = files[0].name
        this.strategyFile = files[0]
      },
      uploadStrategy () {
        this.status = STATUS_UPLOADING
        var formElement = document.querySelector('form')
        var formData = new FormData(formElement)
        var descElement = document.querySelector('#description')
        formData.set('desc', descElement.innerHTML)
        formData.set('author', auth.data.user.getBasicProfile().getName())
        fetch(config.API_BASE_URL + '/upload', {
          method: 'POST',
          body: formData
        }).then((response) => {
          this.uploadResponse = response.json()['message']
          if (response.ok) {
            this.status = STATUS_SUCCESS
          } else {
            this.status = STATUS_FAILED
          }
        }).catch((error) => {
          this.uploadResponse = error
          this.status = STATUS_FAILED
        })
      }
    }
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

  button {
    font-family: 'Avenir', Helvetica, Arial, sans-serif;
    color: #2c3e50;
    height: 40px;
    width: 40%;
    background: none;
    border: 1px solid #2c3e50;
    font-size: 1.5rem;
    cursor: pointer;
  }

  button:hover {
    background: #2c3e50;
    color: #ffffff;
  }

  button:active {
    background: #1e2b3c;
    outline: none;
  }

  .overlay {
    position: fixed;
    top: 0;
    left: 0;
    z-index: 1000;
    height: 100%;
    width: 100%;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .overlay > .dialog {
    display: flex;
    background: white;
    z-index: 1100;
    width: 50%;
    height: 30%;
    min-width: 400px;
    min-height: 200px;
    justify-content: center;
    align-items: center;
  }

  .loader {
    height: 150px;
    width: 150px;
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1200;
  }

  @keyframes rotate {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }
  .loader > .circle {
    width: 140px;
    height: 140px;
    border: 8px solid rgba(0, 0, 0, 0);
    border-radius: 50%;
    border-top: 8px solid #2c3e50;
    z-index: 1400;
    position: absolute;
    animation: rotate 1.5s linear infinite;
    will-change: transform;
  }
  .loader > .inside {
    width: auto;
    z-index: 1300;
    font-size: 1.5rem;
  }

</style>
