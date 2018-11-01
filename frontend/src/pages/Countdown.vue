<template>
    <span>{{timeStamp}}</span>
</template>

<script>
import axios from "axios";
import io from "socket.io-client";

let formatNum = num => {
  if (num < 10) {
    return "0" + num;
  }
  return num;
};

export default {
  props: ["sessionKey"],
  data: function() {
    return {
      session: {},
      secondsLeft: 0,
      connected: false,
      error: null,
      socket: null
    };
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
    axios
      .get(`http://localhost:8081/session/${this.sessionKey}`)
      .then(response => {
        this.session = response.data;
        this.secondsLeft = this.session.secondsLeft;
        this.socket = io("http://localhost:8081");
        this.socket.on("connect", () => {
          console.log("connect");
          this.connected = true;
          this.socket.emit("subscribe", this.session.key);
        });
        this.socket.on("timeUpdate", time => {
          this.secondsLeft = time;
        });
        this.socket.on("disconnect", () => {
          console.log("disconnect");
          this.connected = false;
        });
      });
  }
};
</script>
