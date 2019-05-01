<style scoped lang="less">
.content {
  margin: 3em;
}
.intro-wrapper {
  position: relative;
  .brain {
    background: url(../assets/brain.svg);
  }
  .eye {
    background: url(../assets/eye.svg);
  }
  .brain,
  .eye {
    background-size: contain;
    position: absolute;
    width: 5em;
    height: 5em;
    top: -3em;
  }
  .brain {
    left: -3em;
  }
  .eye {
    right: -3em;
  }
  .intro {
    position: relative;
    z-index: 10;
    border-radius: 0.5em;
    border: 0.15em solid #979797;
    padding: 1em;
    background-color: #fff;
    height: 10em;
    overflow: auto;
  }
}

.btn-primary {
  display: block;
  margin: auto;
}
</style>

<template >
  <div class="welcome-page">
    <div class="content">
      <div class="logo"></div>
      <div class="intro-wrapper">
        <div class="brain"></div>
        <div class="eye"></div>
        <div class="intro">
          <p>
            The contentos harmony puzzle game utilizes fast
            finality of the game to make it frictionless and the
            immutability of ledger to record game state, and
            thus enable trust between players.
          </p>
          <p>Stake. Play. Get rewarded!</p>
        </div>
      </div>
    </div>
    <button class="btn-primary" @click="$emit('join')">Start Game</button>
  </div>
</template>

<script>
export default {
  name: "WelcomePage",

  props: {},
  data() {
    return {};
  },
  mounted: function() {},
  methods: {}
};
</script>
