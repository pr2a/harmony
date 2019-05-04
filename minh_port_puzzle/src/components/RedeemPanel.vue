<style scoped lang="less">
.redeem-panel-container {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
}
.redeem-panel {
  max-width: 100%;
  width: 25em;
  height: 25em;
  background-color: #e2fcf5;
  z-index: 1000;
  border-radius: 0.4em;
  text-align: center;
  color: #1b295e;
  text-transform: uppercase;
  padding: 1em;
  .emphasis,
  .amount {
    font-weight: bold;
  }
  .emphasis {
    font-size: 2em;
  }
  .amount {
    font-size: 3em;
    color: #05b0e9;
  }
  .text {
    p {
      margin: 0.5em 0;
    }
  }
  .email-input {
    border-radius: 0.5em;
    border: 0;
    padding: 1em;
    background-color: #fff;
    display: block;
    width: 100%;
    overflow: auto;
    -webkit-appearance: none;
    font-size: initial;
    outline: none;
    margin-bottom: 1em;
  }
}

::-webkit-input-placeholder {
  text-align: center;
}

:-moz-placeholder {
  /* Firefox 18- */
  text-align: center;
}

::-moz-placeholder {
  /* Firefox 19+ */
  text-align: center;
}

:-ms-input-placeholder {
  text-align: center;
}
</style>

<template >
  <div class="redeem-panel-container flex-hv-center">
    <div class="redeem-panel flex-hv-center">
      <div class="redeem-panel-content">
        <div class="emphasis">You just won</div>
        <div class="amount">{{ reward }}</div>
        <div class="emphasis">Harmony Tokens!</div>
        <div class="text">
          <p>Get your public and private key</p>
          <p>for your claimed token</p>
        </div>
        <input type="text" class="email-input" placeholder="Email..." v-model="email">
        <button class="btn-primary" :disabled="!email.trim()" @click="submitEmail">Submit</button>
      </div>
    </div>
  </div>
</template>

<script>
import store from "../store";
import service from "../service";
export default {
  name: "RedeemPanel",
  props: ["reward"],
  data() {
    return {
      email: "liuzhenglaich@gmail.com"
    };
  },
  methods: {
    submitEmail() {
      service.submitEmail(this.email);
    }
  }
};
</script>
