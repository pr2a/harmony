<style scoped lang="less">
.redeem-panel-container {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  height: 100vh;
  background-color: rgba(0, 0, 0, 0.8);
  z-index: 1;
}
.redeem-panel {
  max-width: 100%;
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
  .amount {
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
    background-color: #fff;
    display: block;
    width: 100%;
    overflow: auto;
    -webkit-appearance: none;
    font-size: initial;
    outline: none;
    margin-bottom: 0.5em;
  }
}

// zien - style message validate email
.err-email {
  font-size: 0.8em;
  color: red;
  padding-bottom: 1em;
  text-transform: none;
  height: 2em;
}
.cancel-email {
  margin: 1em !important;
  text-transform: capitalize;
  color: #482bff;
}

.cancel-email:hover {
  border-bottom: 1px solid #0971f8;
  color: #0971f8;
  cursor: pointer;
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
  <div class="redeem-panel-container flex-hv-center" :style="redeemPanelStyle">
    <div class="redeem-panel flex-hv-center">
      <div class="redeem-panel-content">
        <div class="emphasis" :style="emphasisStyle">You just won</div>
        <div class="amount" :style="amountStyle">{{ reward }}</div>
        <div class="emphasis" :style="emphasisStyle">Harmony Tokens!</div>
        <div class="text" :style="contentEmailStyle">
          <p>Get your public and private key</p>
          <p>for your claimed token</p>
        </div>
        <!-- zien change action @input      @input="validateEmail" v-on:keyup.enter="submitEmail"  -->
        <input
          type="text"
          class="email-input"
          :style="inputEmailStyle"
          placeholder="Email..."
          v-model="email"
          @input="validateEmail"
          v-on:keyup.enter="submitEmail"
        >
        <div class="err-email">{{ err }}</div>

        <!-- zien change   :disabled="(email == '' || err != '')" -->
        <button
          class="btn-primary"
          :disabled="(email == '' || err != '')"
          @click="submitEmail"
          :style="submitButtonStyle"
        >Submit</button>
        <br>
        <br>
        <a class="cancel-email" @click="cancelEmail">Cancel</a>
      </div>
    </div>
  </div>
</template>

<script>
import store from "../store";
import service from "../service";
import { VALIDATE } from "../common/validate";
var dencity = 0;

export default {
  name: "RedeemPanel",
  props: ["reward", "boardSizePx"],
  data() {
    return {
      email: "",
      err: ""
    };
  },
  // mounted: {},
  computed: {
    emphasisStyle() {
      return {
        fontSize: Math.sqrt(window.innerWidth + window.innerHeight) / 2 + "px"
      };
    },
    contentEmailStyle() {
      return {
        fontSize: Math.sqrt(window.innerWidth + window.innerHeight) / 2.5 + "px"
      };
    },
    amountStyle() {
      return {
        fontSize: Math.sqrt(window.innerWidth + window.innerHeight) + "px"
      };
    },
    inputEmailStyle() {
      return {
        padding: Math.sqrt(window.innerWidth + window.innerHeight) / 5 + "px"
      };
    },
    submitButtonStyle() {
      return {
        paddingTop:
          Math.sqrt(window.innerWidth + window.innerHeight) / 6 + "px",
        paddingBottom:
          Math.sqrt(window.innerWidth + window.innerHeight) / 6 + "px",
        paddingLeft:
          Math.sqrt(window.innerWidth + window.innerHeight) / 1 + "px",
        paddingRight:
          Math.sqrt(window.innerWidth + window.innerHeight) / 1 + "px"
      };
    },
    redeemPanelStyle() {
      return { maxWidth: this.boardSizePx };
    }
  },
  methods: {
    submitEmail(e) {
      service.submitEmail(this.email);
    },
    cancelEmail() {
      this.$emit("cancelEmail");
    },
    // zien add action
    validateEmail() {
      if (this.email == "") {
        this.err = "";
      } else {
        let checkEmail = VALIDATE.validateEmail(this.email);

        if (!checkEmail) {
          this.err = "Invalid email";
        } else {
          this.err = "";
        }
      }
    }
  }
};
</script>