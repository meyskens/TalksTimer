<style lang="scss">
@import "../scss/materialize.scss";
</style>

<template>
  <div class="viewer" v-bind:style="{ backgroundColor: color}">
    <div class="container valign-wrapper">
      <div class="row">
          <div class="col s12">
            <p v-if="key" class="countdownView"><countdown v-bind:session-key="key" show-messages=true v-on:color="color = $event"></countdown></p>
          </div>
      </div>
    </div>  
  </div>
</template>

<style lang="scss">
.countdownView {
  font-size: 10em;
}
.viewer {
  height: 100vh;
  margin: 0;
  padding: 0;
}
</style>

<script>
import Countdown from "./Countdown.vue";
export default {
  components: {
    Countdown
  },
  data: function() {
    return {
      key: "",
      color: "",
    }
  },
  computed: {
    realColor: function() {
      if (!this.color) {
        return "#000"
      }
      return this.color
    }
  },
  mounted: function() {
    this.key = this.$route.query.key
  }
};
</script>
