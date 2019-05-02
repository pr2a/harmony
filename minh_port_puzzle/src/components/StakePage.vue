<style scoped lang="less">
.page-content {
  margin: 3em;
}
.msg {
  padding: 0.5em;
  color: #59504d;
  font-family: Fira Sans, sans-serif;
  font-size: 0.8em;
}
.value {
  padding: 1em;
  border-radius: 0.5em;
  color: #59504d;
  font-family: Fira Sans, sans-serif;
  // border: 0.15em solid #979797;
  overflow: auto;
  background-color: #fff;
  margin: 0 auto;
  text-align: center;
}

.btn-primary {
  display: block;
  margin: 0 auto;
}

.action-buttons {
  justify-content: flex-end;
  margin: 0.5em 0;
  .btn-mini {
    margin-left: 0.5em;
  }
}

.multiplier {
  font-size: 1.2em;
}
</style>

<template >
  <div class="tutorial-page">
    <div class="page-content">
      <header class="flex-horizontal">
        <div class="logo"></div>
        <div class="balance info-item">
          <div class="label">Balance</div>
          <div class="content">{{ globalData.balance }}</div>
        </div>
      </header>
      <div class="stake">
        <div class="msg">How many tokens do you want to stake?</div>
        <div class="value">{{ stake }}</div>
        <div class="msg">
          Stake more, win more! You'll get
          <span class="multiplier">{{ stake / 20 }}x</span>
          rewards.
        </div>
        <div class="action-buttons flex-horizontal">
          <button class="btn-mini" @click="minus" :disabled="stake <= 20">-</button>
          <button class="btn-mini" @click="plus" :disabled="stake + 20 > globalData.balance">+</button>
        </div>
      </div>
    </div>
    <button class="btn-primary" @click="stakeToken" :disabled="globalData.balance < 20">Stake</button>
  </div>
</template>

<script>
import store from "../store";
export default {
  name: "TutorialPage",
  data() {
    return {
      globalData: store.data,
      stake: 20
    };
  },
  methods: {
    minus() {
      if (this.stake <= 20) return;
      this.stake -= 20;
    },
    plus() {
      if (this.stake + 20 > this.globalData.balance) return;
      this.stake += 20;
    },
    stakeToken() {
      if (
        confirm(
          `${
            this.stake
          } tokens will be deducted from your balance. Are you sure?`
        )
      ) {
        this.$emit("stake", this.stake);
      }
    }
  }
};
</script>
