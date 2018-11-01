<style lang="scss">
@import "../scss/materialize.scss";
</style>

<template>
  <div>
    <div class="container valign-wrapper">
      <div class="row">
          <div class="col s12">
            <p v-if="key" class="countdownView"><countdown v-bind:session-key="key" show-messages="true"></countdown></p>
          </div>
      </div>
    </div>  
  </div>
</template>

<style lang="scss">
.countdownView {
  font-size: 10em;
  font-family: 'Roboto Mono', monospace;
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
    }
  },
  mounted: function() {
    this.key = this.$route.query.key
  }
};
</script>
