<template>
  <main-layout>
    <div class="row">
      <div class="col s12">
        <div class="card-panel">
          <span>Open the client using: <a v-bind:href="clientURL" target="_new">{{clientURL}}</a></span>
        </div>
      </div>
    </div>
    <div class="row">
      <div class="col s12 m6">
        <div class="card" v-bind:style="{ backgroundColor: timerColor}">
          <div class="card-content">
            <span class="card-title">Timer</span>
            <p v-if="key" class="countdown"><countdown v-bind:session-key="key" v-on:color="timerColor = $event"></countdown></p>
          </div>
        </div>
      </div>
      <div class="col s12 m6">
        <div class="card">
          <div class="card-content">
            <span class="card-title">Control Timer</span>
            <div class="row">
              <div class="col s4">
                <label for="hours">Hours</label>
                <input id="hours" type="number" min="0" v-model="set.hours">
              </div>
              <div class="col s4">
                <label for="minutes">Minutes</label>
                <input id="minutes" type="number" min="0" max="59" v-model="set.minutes">
              </div>
              <div class="col s4">
                <label for="seconds">Seconds</label>
                <input id="seconds" type="number" min="0" max="59" v-model="set.seconds">
              </div>
            </div>
          </div>
          <div class="card-action">
            <a href="#" @click.prevent="resetTimer()">Reset</a>
            <a href="#" @click.prevent="setTimer()">Set</a>
          </div>
        </div>
      </div>
      <div class="col s12 m6">
        <div class="card">
          <div class="card-content">
            <span class="card-title">Show Message</span>
            <div class="row" v-for="(message, index) in messages"  v-bind:key="index">
              <!-- Color Modal -->
              <div v-bind:id="'modal'+index" class="modal color-modal">
                <div class="modal-content">
                  <h4>Set Color</h4>
                  <p><chrome-picker v-model="message.color"/></p>
                </div>
                <div class="modal-footer">
                  <a href="#" class="modal-close waves-effect waves-green btn-flat" @click.prevent="closeModal('modal'+index)">Close</a>
                </div>
              </div>
             
              <div class="col s10">
                <div class="row">
                  <div class="col s05">
                    <a class="btn-color btn-floating waves-effect waves-light" @click.prevent="openModal('modal'+index)" v-bind:style="{'background-color': message.color.hex ? message.color.hex : message.color }"><i class="material-icons black-text">edit</i></a>
                  </div>
                  <div class="col s115">
                    <label for="message">Message</label>
                    <input id="message" type="text" v-model="message.message">
                  </div>
                </div>
              </div>
              <div class="col s2">
                <a class="waves-effect waves-light btn" @click.prevent="sendMessage(message)">Send<i class="material-icons right">send</i></a>
                <a class="btn-floating waves-effect waves-light red message-del" href="#" @click.prevent="deleteMessage(index)"><i class="material-icons">delete</i></a>
              </div>
            </div>
          </div>
          <div class="card-action">
            <a href="#" @click.prevent="addMessage()">Add</a>
            <a href="#" @click.prevent="saveMessages()">Save</a>
            <a href="#" @click.prevent="clearMessage()">Clear</a>
          </div>
        </div>
      </div>
      <div class="col s12 m6">
        <div class="card">
          <div class="card-content">
            <span class="card-title">Colors</span>
            <p v-if="colors.length == 0">No colors set</p>
            <div v-for="(color, index) in colors" v-bind:key="index">
              <div class="row">
                <div class="col s4">
                  <label for="message">Start time</label>
                  <div class="row">
                    <div class="col s6">
                      <label for="minutes">Minutes</label>
                      <input id="minutes" type="number" min="0" v-model="color.minutes">
                    </div>
                    <div class="col s6">
                      <label for="seconds">Seconds</label>
                      <input id="seconds" type="number" min="0" max="59" v-model="color.seconds">
                    </div>
                  </div>
                </div>
                <div class="col s7">
                  <label for="message">Color</label>
                  <slider-picker v-model="color.color" v-bind:swatches="swatches" />
                </div>
                <div class="col s1">
                  <a class="btn-floating waves-effect waves-light red" href="#" @click.prevent="deleteColor(index)"><i class="material-icons">delete</i></a>
                </div>
              </div>
            </div>
          </div>
          <div class="card-action">
            <a href="#" @click.prevent="addColor()">Add</a>
            <a href="#" @click.prevent="saveColors()">Save</a>
          </div>
        </div>
      </div>
    </div>    
  </main-layout>
</template>

<style lang="scss">
.countdown {
  font-size: 4em;
  text-align: center;
}
.vc-slider {
  width: 100%;
}
.message-del {
  margin-left: 10px;
}
.color-modal {
  width: 280px !important;
}
.row .col.s05 {
    width: 5%;
    margin-left: auto;
    left: auto;
    right: auto;
}
.row .col.s115 {
    width: 95%;
    margin-left: auto;
    left: auto;
    right: auto;
}
</style>

<script>
import { Slider, Chrome } from 'vue-color'
import MainLayout from "../layouts/Main.vue";
import Countdown from "./Countdown.vue";
import * as Timer from "../api/timer.js"
import * as Messages from "../api/messages.js"
import * as Colors from "../api/colors.js"

export default {
  components: {
    MainLayout,
    Countdown,
    'slider-picker': Slider,
    'chrome-picker': Chrome,
  },
  computed: {
    clientURL: function() {
      return `${window.location.origin}/#/viewer/?key=${this.key}`
    },
  },
  data: function() {
    return {
      key: "",
      set: {
        hours: 0,
        minutes: 0,
        seconds: 0,
      },
      colors: [],
      messages: [],
      swatches: ['.90', '.65', '.50', '.35', '.20'],
      timerColor: "",
    }
  },
  mounted: function() {
    try {
      window.$ = window.jQuery = require('jquery');
      window.materialize = require('materialize-css');
    } catch (e) {
        console.log(e);
    }

    this.key = this.$route.query.key
    Colors.get(this.key).then(response => {
      this.colors = response.data.options || []
      for (let color of this.colors) {
        color.seconds = color.from % 60
        color.minutes = Math.floor(color.from / 60)
      }
    })
    Messages.get(this.key).then(response => {
      this.messages = response.data.messages || []
    })
  },
  methods: {
    setTimer: function() {
      let time = parseInt(this.set.seconds)
      console.log(this.set)
      time += parseInt(this.set.minutes) * 60
      time += parseInt(this.set.hours) * 60 * 60

      Timer.setTimer(this.key, time)
    },
    resetTimer : function() {
      Timer.setTimer(this.key, 0)
    },
    sendMessage : function(message) {
      if (message.color.hex) {
        message.color = message.color.hex
      }
      Messages.send(this.key, message)
    },
    clearMessage : function() {
      Messages.send(this.key, "")
    },
    addMessage: function() {
      this.messages.push(
        {
          message: "",
          color: "#ffffff",
        }
      )
    },
    deleteMessage: function(index) {
      let partOne = this.messages.slice(0, index);
      let partTwo = this.messages.slice(index+1, this.messages.length);
      this.messages = partOne.concat(partTwo)
    },
    saveMessages: function() {
      for (let message of this.messages) {
        if (message.color.hex) {
          message.color = message.color.hex
        }
      }
      Messages.set(this.key, this.messages).then(() => {
        M.toast({html: 'Messages saved'})
      })
    },
    addColor: function() {
      this.colors.push({
        color: "#fff",
        from: 0,
        seconds: 0,
        minutes: 0,
      })
    },
    deleteColor: function(index) {
      let partOne = this.colors.slice(0, index);
      let partTwo = this.colors.slice(index+1, this.colors.length);
      this.colors = partOne.concat(partTwo)
    },
    saveColors: function() {
      for (let color of this.colors) {
        if (color.color.hex) {
          color.color = color.color.hex
        }
        color.from = parseInt(color.seconds) + (parseInt(color.minutes) * 60)
      }
      Colors.set(this.key, this.colors).then(() => {
        M.toast({html: 'Colors saved'})
      })
    },
    openModal: function(index) {
      $(`#${index}`).modal()
      $(`#${index}`).modal("open")
    },
    closeModal: function(index) {
      $(`#${index}`).modal("close")
    },
  }
};
</script>
