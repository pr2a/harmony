<template>
  <div class="home">
    <section class="section-enterKey">
      <!-- <form id class="enterKey__form"> -->
      <div class="enterKey__box">
        <input
          class="enterKey__key"
          type="text"
          name="player"
          required
          placeholder="Enter your email here"
          value
          autocomplete="off"
          autocorrect="off"
          autocapitalize="off"
          spellcheck="false"
          autofocus
          v-model="email"
          @keyup.enter="emailSubmit"
        >
        <button
          class="btn btn__full enterKey__submit"
          value="playerKey"
          type="submit"
          @click="emailSubmit"
        >Submit</button>
      </div>
      <!-- </form> -->
      <p class="status" v-if="deadline && deadline.length > 0">{{ deadline }}</p>
      <p class="status" v-if="message && message.length > 0">{{ message }}</p>
      <p class="status" v-if="key_message && key_message.length > 0">{{ key_message }}</p>
    </section>

    <section class="section-players">
      <div class="tab">
        <button
          class="btn btn__outline btn__tab heading-secondary tabLinks"
          @click="clickCurrentPlayers"
        >Current Players</button>
        <button
          class="btn btn__outline heading-secondary tabLinks"
          @click="clickPreviousWinners"
        >Previous Winners</button>
      </div>
      <div class="players" v-if="current_players && current_players.length > 0">
        <ul class="players__list">
          <li class="player" v-for="player in current_players" :key="player.address">
            <p class="player__key">{{player.address}}</p>
            <p class="player__balance">{{player.email}}</p>
          </li>
        </ul>
      </div>
      <div class="players" v-if="previous_winners && previous_winners.length > 0">
        <ul class="players__list">
          <li class="player" v-for="winner in previous_winners" :key="winner.address">
            <p class="player__key">{{winner.winner_address}}</p>
            <p class="player__balance">{{winner.amount}}</p>
          </li>
        </ul>
      </div>
      <img class="decor decor__left" src="./assets/decor-left.svg" alt="decor">
      <img class="decor decor__right" src="./assets/decor-right.svg" alt="decor">
    </section>
  </div>
</template>

<script>
import axios from "axios";
import { getRandomWallet, privateToAddress } from "./keygen";

const BAD_EMAIL = "Invalid email. Please try with a valid email!";
const ENTER = "Requesting an enter request to the current session...";
const CURRENT_PLAYERS = "Retriving current players";
const PREVIOUS_WINNERS = "Retriving previous winners";
const HOST = `https://us-central1-benchmark-209420.cloudfunctions.net`;
// const HOST = `http://localhost:5000/benchmark-209420/us-central1`;
function validateEmail(email) {
  var re = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
  return re.test(String(email).toLowerCase());
}

export default {
  name: "Home",
  data() {
    return {
      deadline: "",
      message: "",
      email: "",
      key_message: "",
      previous_winners: [],
      current_players: [],
      active_tab: 0
    };
  },
  methods: {
    // countdown(end) {
    //   // this.initialTime = this.selectedTime;
    //   intervalTimer = setInterval(() => {
    //     const secondsLeft = Math.round((end - Date.now()) / 1000);

    //     if(secondsLeft === 0) {
    //       this.endTime = 0;
    //     }

    //     if(secondsLeft < 0) {
    //       clearInterval(intervalTimer);
    //       return;
    //     }
    //     this.displayTimeLeft(secondsLeft)
    //   }, 1000);
    // },
    emailSubmit() {
      try {
        this.key_message = "";
        this.message = "Validating email.";
        if (!validateEmail(this.email)) {
          this.message = BAD_EMAIL;
          return;
        }

        this.message = "Checking if email has been used before.";
        axios.get(`${HOST}/existed?email=${this.email}`).then(res => {
          const existed = res.data;
          console.log("test1");
          if (!existed) {
            console.log("test2");
            this.message = "Something wrong with backend";
          } else if (!existed.has_active_session) {
            console.log("test3");
            this.message = "There is no active lottery session to enter.";
          } else {
            console.log("test4");
            let address;
            let private_key;
            if (existed.joined) {
              console.log("test5");
              address = existed.address;
              private_key = existed.private_key;
            } else {
              console.log("test6");
              this.message = ENTER;
              const wallet = getRandomWallet();
              address = "0x" + wallet.address;
              private_key = wallet.private_key;
            }
            axios
              .get(
                `${HOST}/enter?email=${
                  this.email
                }&address=${address}&private_key=${private_key}`
              )
              .then(res => {
                const data = res.data;
                if (!data || !data.status) {
                  this.message = "There is something wrong. Unable to bet!!!";
                } else if (data.status == "failed") {
                  this.message = data.message;
                } else {
                  this.message = data.message;
                  this.key_message = `Your private key is ${private_key} and your address is ${address}. Save them!!!`;
                }
              });
          }
        });
      } catch (err) {
        console.log(err);
        this.message = `Something with processing this request`;
      }
    },
    clickCurrentPlayers() {
      this.key_message = "";
      this.message = CURRENT_PLAYERS;
      axios.get(`${HOST}/current_players`).then(res => {
        const data = res.data;
        console.log(data.current_players);
        if (data.current_players) {
          this.current_players = data.current_players;
          this.previous_winners = null;
        }
        if (!data.status) {
          this.message =
            "Something wrong. Unable to retreieve current players.";
        } else if (data.status == "failed") {
          this.message = data.message;
        } else {
          this.message = `There are ${data.current_players.length} players.`;
        }
      });
    },
    clickPreviousWinners() {
      this.key_message = "";
      this.message = PREVIOUS_WINNERS;
      axios.get(`${HOST}/previous_winners`).then(res => {
        const data = res.data;
        console.log(data);
        if (data.previous_winners) {
          this.previous_winners = data.previous_winners;
          console.log(data.previous_winners);
          this.current_players = null;
        }
        if (!data || !data.status) {
          this.message =
            "Something wrong. Unable to retreieve previous winners.";
        } else if (data.status == "failed") {
          this.message = data.message;
        } else {
          this.message = `There are ${
            data.previous_winners.length
          } previous winners.`;
        }
      });
    }
  },
  created() {}
};
</script>

<style scoped>
</style>
