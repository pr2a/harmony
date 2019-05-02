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
    <stake-page @stake="stake" v-if="step === 3"></stake-page>
    <tutorial-page @start="startGame" v-if="step === 4"></tutorial-page>
    <puzzle-page @restart="restartGame" v-if="step === 5"></puzzle-page>
  </div>
</template>

<script>
import WelcomePage from "./WelcomePage";
import PuzzlePage from "./PuzzlePage";
import EmailPage from "./EmailPage";
import KeyPage from "./KeyPage";
import TutorialPage from "./TutorialPage";
import StakePage from "./StakePage";
import service from "../service";

export default {
  name: "HostingPage",
  components: {
    WelcomePage,
    EmailPage,
    KeyPage,
    TutorialPage,
    StakePage,
    PuzzlePage
  },
  data() {
    return {
      step: 3,
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
    stake(value) {
      service.stakeToken(value).then(() => {
        if (localStorage.getItem("hideTutorial")) this.step++;
        this.step++;
      });
    },
    startGame() {
      this.step++;
    },
    restartGame() {
      this.step = 3;
    }
  }
};
</script>
