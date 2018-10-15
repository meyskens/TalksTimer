<template>
  <main-layout>
    <p>{{this.$route.query.key}}</p>
  </main-layout>
</template>

<script>
import MainLayout from "../layouts/Main.vue";

export default {
  components: {
    MainLayout
  }
};
</script>
