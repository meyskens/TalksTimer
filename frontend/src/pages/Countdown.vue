<template>
    <span><span v-if="!message" class="numbers">{{timeStamp}}</span><span v-if="message">{{message}}</span></span>
</template>

<style lang="scss">
.numbers {
  font-family: 'Roboto Mono', monospace;
}
</style>


<script>
import io from "socket.io-client";
import { URL } from "../api/const.js"
import * as Timer from "../api/timer.js";

let formatNum = num => {
  if (num < 10) {
    return "0" + num;
  }
  return num;
};

export default {
  props: {
    sessionKey: String,
    showMessages: Boolean,
  },
  data: function() {
    return {
      session: {},
      secondsLeft: 0,
      connected: false,
      error: null,
      socket: null,
      message: "",
      messageClearTimer: null,
    };
  },
  methods: {
    clearMessage: function() {
      this.message = ""
    }
  },
  computed: {
    timeStamp: function() {
      let seconds = parseInt(this.secondsLeft, 10); // copies
      let hrs = Math.floor(seconds / 3600);
      seconds -= hrs * 3600;
      let mnts = Math.floor(seconds / 60);
      seconds -= mnts * 60;
      return `${hrs > 0 ? formatNum(hrs) + ":" : ""}${formatNum(
        mnts
      )}:${formatNum(seconds)}`;
    }
  },
  mounted: function() {
    Timer.get(this.sessionKey).then(response => {
        this.session = response.data;
        this.secondsLeft = this.session.secondsLeft;
        this.socket = io(URL);
        this.socket.on("connect", () => {
          console.log("connect");
          this.connected = true;
          this.socket.emit("subscribe", this.session.key);
        });
        this.socket.on("timeUpdate", time => {
          this.secondsLeft = time;
        });
        if (this.showMessages) {
          this.socket.on("message", message => {
            this.message = message;
            clearTimeout(this.messageClearTimer)
            this.messageClearTimer = setTimeout(this.clearMessage, 10000)
          });
        }
        this.socket.on("disconnect", () => {
          console.log("disconnect");
          this.connected = false;
        });
      });
  }
};
</script>
