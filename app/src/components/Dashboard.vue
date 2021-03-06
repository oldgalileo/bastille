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
        <div class="button" @click="reset()">Close</div>
      </div>
      <div class="content" v-else-if="isFailed">
        <h2>Uh-oh!</h2>
        <p>Could not upload strategy...</p>
        <div class="button" @click="reset()">Close</div>
      </div>
    </div>
  </div>
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
            <input type="file" class="file-input" id="exec" name="exec" hidden="true" @change="handleFile( $event.target.files )">
            <template v-if="strategyFileName === ''">
              <label for="exec" id="exec-button" class="file-label">Select File...</label>
            </template>
            <template v-else>
              <label for="exec" id="exec-button" class="file-label">{{ strategyFileName }}</label>
            </template>
          </li>
          <li class="input-line">
            <div class="submit" @click="uploadStrategy()"><div class="inner"></div>Submit</div>
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
  .form {
    max-width: 700px;
    width: 80%;
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
    flex-flow: column nowrap;
    justify-content: flex-start;
    align-items: flex-start;
  }

  .header-line {
    justify-content: center;
  }

  .header-line .header {
    font-size: 2.0rem;
    font-weight: bold;
    padding-bottom: 15px;
  }

  .input-line label:not(.header) {
    font-size: 0.75rem;
    line-height: 1.5rem;
    text-align: left;
  }

  .input-line input[type=text] {
    font-family: 'Avenir', Helvetica, Arial, sans-serif;
    border: 2px solid #2c3e50;
    display: flex;
    justify-content: center;
    align-items: center;
    -webkit-appearance: none;
    -moz-appearance: none;
    width: 100%;
    padding: 15px;
    padding-top: 15px;
    padding-bottom: 15px;
    border-radius: 3px;
    font-size: 1rem;
    margin-bottom: 10px;
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
    padding: 15px;
    border: 2px solid #2c3e50;
    padding-top: 15px;
    padding-bottom: 15px;
    border-radius: 3px;
    font-size: 1.0rem;
    text-align: left;
    font-weight: 400;
    margin-bottom: 10px;
    overflow: auto;
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
    width: 100%!important;
    border: none!important;
    margin-bottom: 10px;
    display: flex;
    align-items: center;
    font-size: 1rem!important;
    padding: 15px;
    padding-left: 0;
    padding-top: 10px;
    padding-bottom: 10px;
    font-weight: bold;
    cursor: pointer;
  }

  .input-line:last-child {
    justify-content: center;
  }

  .submit {
    font-family: 'Avenir', Helvetica, Arial, sans-serif;
    color: #2c3e50;
    padding: 15px;
    padding-top: 10px;
    padding-bottom: 10px;
    background: none;
    border: 2px solid #2c3e50;
    font-weight: bold;
    border-radius: 3px;
    font-size: 1rem;
    cursor: pointer;
    position: relative;
    user-select: none;
  }
  .submit .inner {
    left: 0;
    top: 0;
    background-color: #2c3e50;
    width: 0%;
    overflow: hidden;
    display: flex;
    box-sizing: border-box;
    position: absolute;
    content: ' ';
    color: white;
    transition: width 0.2s ease,
                background-color 0.2s ease;
    will-change: width;
    z-index: 1;
    color: white;
    content: 'Submit';
  }
  .submit .inner::after {
    color: white;
    content: 'Submit';
    display: block;
    padding: 15px;
    padding-top: 10px;
    padding-bottom: 10px;
  }
  .submit:hover .inner {
    width: 100%;
  }
  .submit:active .inner {
    background-color: #1e2b3c;
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
