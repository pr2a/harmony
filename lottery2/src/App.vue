<template>
  <div id="app">
    <Header/>
    <Home/>
  </div>
</template>

<script>
import Header from "./Header";
import Home from "./Home";

export default {
  name: "app",
  components: {
    Header,
    Home
  },
  data() {
    return {
      msg: "Welcome to Your Vue.js App"
    };
  }
};
</script>

<style>
</style>
