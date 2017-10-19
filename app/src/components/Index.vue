<template>
  <div class="main">
    <div class="overlay" v-if="popup===true" @click="toggleOverlay">
      <div class="dialog">
        <div class="content" v-if="uploading==='RESPONSE'">
          <div class="loader">
            <div class="inside">Uploading</div>
            <div class="circle"></div>
          </div>
        </div>
        <div class="content" v-else>
          <h2>Success!</h2>
          <p>Your strategy was uploaded successful.</p>
        </div>
      </div>
    </div>
    <section>
      <div class="content">
        <div class="content">
          <h1 class="big-header">Bastille</h1>
          <p>A new approach to programmatic IPD (Iterative Prisoner's Dilemma)</p>
        </div>
        <!--<div class="scroll">V</div>-->
      </div>
    </section>
    <!--<section class="gray">-->
      <!--<div class="split left">-->
        <!--<h1 class="header">Submit</h1>-->
        <!--<div class="content">-->
          <!--<p class="paragraph" style="font-size: 1.5rem; font-weight: bold;">Instructions –</p>-->
          <!--<p class="paragraph">Please ensure that your code is in an executable format. To test if it can be executed, open up "Terminal" and run:</p>-->
          <!--<p class="paragraph" style="text-indent: 50px;">sh -c "/path/to/your/strategy"</p>-->
          <!--<p class="paragraph">If this command succeeds, you can upload your code. Just click the button below and select your file, kick back, and watch yourself lose to always defect.</p>-->
          <!--<form enctype="multipart/form-data">-->
            <!--<input type="file" name="file" id="strategy" @change="handleUpload( $event.target.files )" class="strategy-input">-->
            <!--<label for="strategy" class="button">Submit</label>-->
          <!--</form>-->
        <!--</div>-->
      <!--</div>-->
      <!--<span class="divider"></span>-->
      <!--<div class="split right">-->
        <!--<h1 class="header">Leaderboard</h1>-->
        <!--<div class="content">-->
          <!--<h3>Coming Soon...</h3>-->
        <!--</div>-->
      <!--</div>-->
    <!--</section>-->
  </div>
</template>

<script>
import * as config from '../config'

export default {
  name: 'Index',
  data () {
    return {
      popup: false,
      uploading: 'INACTIVE',
      uploadingResponse: '',
      quotes: [
        'Patrick is already disappointed...',
        'In everything one thing is impossible: rationality.',
        'Inconceivable!',
        '"You only think I guessed wrong! That\'s what\'s so funny!" - Vizzini'
      ]
    }
  },
  methods: {
    reset: function () {
      this.uploading = 'INACTIVE'
      this.uploadingResponse = ''
      this.popup = false
    },
    handleUpload: function (fileList) {
      console.log('Halp')
      this.uploading = 'UPLOADING'
      this.toggleOverlay()
      this.uploadFile(fileList[0])
    },
    uploadFile: function (file) {
      const formData = new FormData()
      formData.append('strategy', file, file.name)
      fetch(config.API_BASE_URL + '/upload', {
        method: 'POST',
        body: formData
      }).then(
        response => {
          console.log(response)
          this.uploadingResponse = response.json()['message']
        }
       ).then(
        this.uploading = 'RESPONSE'
      ).catch(
        error => console.log(error)
      )
    },
    toggleOverlay: function () {
      console.log('Running')
      if ((this.uploading !== 'UPLOADING') && this.popup) {
        this.popup = false
      } else {
        this.popup = true
      }
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.divider {
  height: 25%;
  width: 1px;
  background: #2c3e50;
  align-self: center;
}

.button {
  align-content: center;
  text-align: center;
  border: thin solid #2c3e50;
  min-width: 60px;
  line-height: 30px;
  min-height: 30px;
  padding: 4px 10px;
  font-size: 1.25rem;
  cursor: pointer;
}

.button:hover {
  background: #2c3e50;
  color: white;
}

.button:active {
  background: #202e41;
}

.main {
  height: 100%;
  width: 100%;
  display: flex;
  flex-flow: row wrap;
  justify-content: space-between;
  align-items: center;
}

.scroll {
  font-weight: bold;
  line-height: 20px;
  transform:scale(4,1.5);
  opacity: 0.75;
  padding-bottom: 20px;
}

section {
  width: 100%;
  height: 100%;
  display: flex;
  flex-flow: row wrap;
  justify-content: center;
}

.content {
  width: 100%;
  height: 100%;
  display: flex;
  flex-flow: column nowrap;
  align-items: center;
  justify-content: center;
}

.content > .paragraph {
  width: 75%;
  text-align: justify;
  margin-top: 0;
}

section .big-header {
  font-size: 5rem;
  margin-bottom: 0;
}

section.gray {
  background: #2c3e50;
  color: white;
}

section.gray .button {
  border: thin solid white;
}

section.gray .button:hover {
  background: white;
  color: #2c3e50;
}

section.gray .button:active {
  background: #f0f0f0;
  border-color: #f0f0f0;
  color: #2c3e50;
}

section.gray .divider {
  background: white;
}

.split {
  display: flex;
  flex-flow: column nowrap;
  width: 39%;
  height: 50%;
  justify-content: center;
  align-items: center;
  align-self: center;
  margin: 0;
}

.split > .header {
  font-size: 2.5rem;
  margin-top: 50px;
  margin-bottom: auto;
}

.left {
  margin-left: 10%;
}
.right {
  margin-right: 10%;
}

.strategy-input {
  width: 0.1px;
  height: 0.1px;
  opacity: 0;
  overflow: hidden;
  position: absolute;
  z-index: -1;
}

.strategy-input:focus + label, .strategy-input + label:hover {

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
