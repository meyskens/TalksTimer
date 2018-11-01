<template>
  <main-layout>
    <div class="container">
      <h2 class="center-align">Welcome to TalkTimer</h2>
      <p class="center-align">TalkTimer is a remote controllable timer designed to signal the time and messages to presenters</p>
      <p class="center-align">
        <a class="waves-effect waves-light btn-large"  v-on:click="newSession"><i class="material-icons left">add_circle</i>Start Session</a>
      </p>
    </div>
  </main-layout>
</template>

<script>
import axios from "axios";
import MainLayout from "../layouts/Main.vue";
import routes from "../routes";
import * as Timer from "../api/timer.js";

export default {
  components: {
    MainLayout
  },
  methods: {
    newSession: function() {
      Timer.create().then(response => {
        // handle success
        if (response.data.instance === "dev-server") {
          this.$router.push({
            path: "controller",
            query: { key: response.data.key }
          });
        } else {
          window.location.replace(`https://${response.data.instance}.talkstimer.com/#/controller?key=${response.data.key}`);
        }
      });
    }
  }
};
</script>
