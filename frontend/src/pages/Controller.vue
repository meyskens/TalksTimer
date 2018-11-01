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
        <div class="card">
          <div class="card-content">
            <span class="card-title">Timer</span>
            <p v-if="key" class="countdown"><countdown v-bind:session-key="key"></countdown></p>
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
            <div>
              <label for="message">Message</label>
              <input id="message" type="text" v-model="messageField">
            </div>
          </div>
          <div class="card-action">
            <a href="#" @click.prevent="clearMessage()">Clear</a>
            <a href="#" @click.prevent="sendMessage('Repeat the quesion!')">Repeat Question</a>
            <a href="#" @click.prevent="sendMessage(messageField)">Send</a>
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
</style>

<script>
import MainLayout from "../layouts/Main.vue";
import Countdown from "./Countdown.vue";
import * as Timer from "../api/timer.js"
import * as Messages from "../api/messages.js"

export default {
  components: {
    MainLayout,
    Countdown,
  },
  computed: {
    clientURL: function() {
      return `${window.location.origin}/#/viewer/?key=${this.key}`
    }
  },
  data: function() {
    return {
      key: "",
      set: {
        hours: 0,
        minutes: 0,
        seconds: 0,
      },
      messageField: "",
    }
  },
  mounted: function() {
    this.key = this.$route.query.key
  },
  methods: {
    setTimer: function() {
      let time = this.set.seconds
      time += this.set.minutes * 60
      time += this.set.hours * 60 * 60

      Timer.setTimer(this.key, time)
    },
    resetTimer : function() {
      Timer.setTimer(this.key, 0)
    },
    sendMessage : function(message) {
      Messages.send(this.key, message)
    },
    clearMessage : function() {
      Messages.send(this.key, "")
    },
  }
};
</script>
