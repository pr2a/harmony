<style scoped lang="less">
.host {
  max-width: 800px;
  margin: 0 auto;
}
</style>

<template >
  <div class="host">
    <welcome-page @join="join" v-if="step === 0"></welcome-page>
    <email-page @submit="submitEmail" v-if="step === 1"></email-page>
    <key-page :userKey="userKey" v-if="step === 2" @start="startGame"></key-page>
    <tutorial-page @start="startGame" v-if="step === 3"></tutorial-page>
    <puzzle-page v-if="step === 4"></puzzle-page>
  </div>
</template>

<script>
import WelcomePage from "./WelcomePage";
import PuzzlePage from "./PuzzlePage";
import EmailPage from "./EmailPage";
import KeyPage from "./KeyPage";
import TutorialPage from "./TutorialPage";
import service from "../service";

export default {
  name: "HostingPage",
  components: {
    WelcomePage,
    EmailPage,
    KeyPage,
    TutorialPage,
    PuzzlePage
  },
  data() {
    return {
      step: 0,
      userKey: "Oxhsa89sd23jkl3450stypose00"
    };
  },
  mounted: function() {},
  methods: {
    join() {
      this.step++;
    },
    submitEmail(email) {
      service.register(email).then(() => {
        this.step++;
      });
    },
    startGame() {
      this.step++;
    }
  }
};
</script>
