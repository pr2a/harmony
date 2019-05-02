<style scoped lang="less">
.page-content {
  margin: 3em;
}

.logo {
  font-size: 2em;
}

@media (min-width: 800px) {
  .logo {
    font-size: 2.5em;
  }
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
footer {
  position: fixed;
  bottom: 0;
  left: 0;
  width: 100%;
  padding: 2em 1em;
  .link {
    font-size: 0.8em;
    text-align: center;
  }
}
.host {
  max-width: 600px;
  margin: 0 auto;
}
.info-item > .content {
  font-size: 1.4em;
}
.btn-mini.start-btn {
  width: auto;
  padding: 0 1em;
}
.stake-amount {
  margin: 0 0.5em;
  background-color: #69fabd;
  border-radius: 0.3em;
  border: 0;
  color: #19586d;
  height: 2em;
  width: 5em;
  font-size: 0.8em;
}
.stake-row {
  margin: 1em auto 0;
}
</style>

<template >
  <div class="flex-horizontal stake-row">
    <div class="action-buttons flex-horizontal flex-grow">
      <button class="btn-mini" @click="minus" :disabled="stake <= 20">
        <font-awesome-icon icon="minus"></font-awesome-icon>
      </button>
      <div class="stake-amount flex-hv-center">{{ stake }}</div>
      <button class="btn-mini" @click="plus" :disabled="stake + 20 > globalData.balance">
        <font-awesome-icon icon="plus"></font-awesome-icon>
      </button>
    </div>
    <button class="btn-mini start-btn" @click="stakeToken" :disabled="globalData.balance < 20">Start</button>
  </div>
</template>

<script>
import service from "../service";
import store from "../store";
export default {
  name: "StakeRow",
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
      service.stakeToken(this.stake).then(() => {
        this.$emit("stake", this.stake);
      });
    }
  }
};
</script>
