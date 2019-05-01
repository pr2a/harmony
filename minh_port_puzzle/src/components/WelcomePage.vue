<style scoped lang="less">
.content {
  margin: 3em;
}
.intro {
  border-radius: 0.5em;
  border: 0.15em solid #979797;
  padding: 1em;
  background-color: #fff;
  height: 10em;
  overflow: auto;
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
      <div class="intro">
        The contentos harmony puzzle game utilizes fast
        finality of the game to make it frictionless and the
        immutability of ledger to record game state, and
        thus enable trust between players.
        "In order to play the game the player has to stake
        tokens. Harmony would be able to confirm the
        staking transaction instantly. The player can then
        proceed to be matched with another player. The
        other player is also ready to play instantly as their
        staking transaction has also gone through. Because The contentos harmony puzzle game utilizes fast
        finality of the game to make it frictionless and the
        immutability of ledger to record game state, and
        thus enable trust between players.
        "In order to play the game the player has to stake
        tokens. Harmony would be able to confirm the
        staking transaction instantly. The player can then
        proceed to be matched with another player. The
        other player is also ready to play instantly as their
        staking transaction has also gone through. Because The contentos harmony puzzle game utilizes fast
        finality of the game to make it frictionless and the
        immutability of ledger to record game state, and
        thus enable trust between players.
        "In order to play the game the player has to stake
        tokens. Harmony would be able to confirm the
        staking transaction instantly. The player can then
        proceed to be matched with another player. The
        other player is also ready to play instantly as their
        staking transaction has also gone through. Because
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
