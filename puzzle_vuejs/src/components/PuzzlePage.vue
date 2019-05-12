<style scoped lang="less">
.score-container {
  margin-right: auto;
  margin-left: auto;
  justify-content: space-between;
  display: flex;
  flex: 1;
  flex-direction: row;
  min-height: 55px;
}

footer {
  margin: 1em auto 0;
  .btn-primary {
    font-size: 1em;
    background-color: #482bff;
  }
}

.board-wrapper {
  position: relative;
  margin: 0 auto;
  flex-grow: 0;
  flex-shrink: 0;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s;
}
.fade-enter, .fade-leave-to /* .fade-leave-active below version 2.1.8 */ {
  opacity: 0;
}

.game-over-message {
  font-weight: bold;
  text-align: center;
  background-color: rgba(255, 255, 255, 0.7);
  border-radius: 0.3em;
}

.content-tutorial {
  height: 100%;
  padding: 0 12px;
  background: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
}

.main-container {
  height: 100vh;
  .game-container {
    position: relative;
    display: flex;
    flex-direction: column;
    height: 100vh;
    .overlay {
      width: 100%;
      height: 100%;
      position: absolute;
      z-index: 2;
      display: flex;
      align-items: center;
      justify-content: center;
    }
  }
}

.blur-text {
  opacity: 0.8;
}

.appearing {
  animation: appearing 1s;
  -webkit-animation: appearing 1s;
}

@keyframes appearing {
  0% {
    opacity: 0;
  }

  100% {
    opacity: 1;
  }
}

.action-row + .action-row {
  margin-top: 1em;
}
.info-item {
  align-self: flex-end;
  .content {
    position: relative;
  }
}
.count-down {
  .seconds-left {
    &.hurry-up {
      color: #f6371d;
    }
    &.hurry-up {
      animation-name: headShake;
      animation-duration: 1s;
      animation-timing-function: ease-int-out;
      animation-iteration-count: infinite;
    }
  }
}

@keyframes headShake {
  0% {
    transform: translateX(0);
  }

  6.5% {
    transform: translateX(-6px) rotateY(-9deg);
  }

  18.5% {
    transform: translateX(5px) rotateY(7deg);
  }

  31.5% {
    transform: translateX(-3px) rotateY(-5deg);
  }

  43.5% {
    transform: translateX(2px) rotateY(3deg);
  }

  50% {
    transform: translateX(0);
  }
}

.number-increase {
  display: block;
  position: absolute;
  left: 0;
  top: 0;
  color: #2c3e50;
  width: 100%;
  animation: up-disappear 1.5s;
}
@keyframes up-disappear {
  0% {
    opacity: 0.7;
  }

  100% {
    opacity: 0;
    transform: translateY(-40px);
  }
}
.link-footer {
  left: 0;
  width: 100%;
  padding: 1em;
  text-align: center;
  z-index: 1000;
  flex: 1;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  min-height: 40px;
}
.fake-footer {
  flex: 1;
}
.icon-clock,
.icon-token {
  background-size: contain;
}
.icon-clock {
  background-image: url(../assets/clock.svg);
}
.icon-token {
  background-image: url(../assets/token.svg);
}
.level-text {
  font-weight: bold;
}
.logo {
  align-self: flex-start;
}
.link {
  font-size: 0.8em;
  text-align: center;
  text-decoration: none;
}
</style>

<template>
  <div id="app">
    <redeem-panel
      v-if="gameEnded && !globalData.email && !cancelEmail"
      :reward="reward"
      :boardSizePx="boardSizePx"
      @cancelEmail="closeEmailPopup"
    ></redeem-panel>
    <div class="main-container appearing">
      <div class="game-container" ref="gameContainer">
        <!--<redeem-panel v-if="gameEnded && !globalData.email" :reward="reward"></redeem-panel>-->
        <a
          :href="'https://explorer.harmony.one/#/address/' + globalData.address"
          class="logo"
          target="_blank"
        ></a>
        <div class="score-container" :style="{ width: boardSizePx + 'px' }">
          <div class="balance info-item">
            <div class="label">
              <div class="icon-token" :style="iconTokenStyle"></div>
            </div>
            <div class="content">
              {{ globalData.balance }}
              <transition>
                <span v-if="balanceIncrease!=''" class="number-increase">
                  {{
                  balanceIncrease
                  }}
                </span>
              </transition>
            </div>
          </div>
          <a
            :style="titleStyle"
            :href="'https://explorer.harmony.one/#/address/' + globalData.address"
            class="logo"
            target="_blank"
          ></a>
          <div class="count-down info-item" :style="infoItemStyle">
            <div class="label">
              <div class="icon-clock" :style="iconTokenStyle"></div>
            </div>
            <div class="content">
              <div
                class="seconds-left"
                :class="{ 'hurry-up': secondsLeft && secondsLeft <= 12, 'game-over': !secondsLeft }"
              >{{ secondsLeft | time }}</div>
              <transition>
                <span v-if="timeIncrease!=''" class="number-increase">
                  {{
                  timeIncrease
                  }}
                </span>
              </transition>
            </div>
          </div>
        </div>

        <div class="board-wrapper" :style="boardWrapperStyle">
          <div v-if="gameEnded || !gameStarted">
            <div class="overlay game-over-message appearing">
              <div class="content content-tutorial">
                <p :style="gameOverStyle" v-if="!globalData.account">Logging in...</p>
                <p :style="gameOverStyle" v-else-if="gameEnded">Game over!</p>
                <p class="blur-text" :style="gameTutorialStyle" v-else-if="!gameStarted">
                  <span :style="gameTutorialSmallStyle"
                  >Move cursor to adjacent cells to increase the number by 1. Win a level by making all numbers equal!</span>
                  <br>
                  <br>Place bet (bottom left) and click “Start"
                </p>
              </div>
            </div>
          </div>
          <transition name="fade" v-for="(level, i) in levels" :key="i">
            <Game
              :ref="'game' + i"
              :listen-own-key-events-only="false"
              :tab-index="1"
              :board-size-px="boardSizePx"
              :game="level"
              :gameLevel="levelIndex+1"
              :gameStarted="gameStarted"
              :gameEnded="gameEnded"
              @completeLevel="onLevelComplete"
              v-if="i === levelIndex"
            ></Game>
          </transition>
        </div>
        <stake-row v-if="!gameStarted" @stake="startGame" :style="stakeRowStyle" @stakeToken="resetLevel"></stake-row>
        <footer class="flex-vertical" :style="{ width: boardSizePx + 'px' }" v-if="gameStarted">
          <div class="flex-horizontal action-row">
            <span
              class="flex-grow level-text"
              :style="levelTextStyle"
            >Level: {{ levelIndex + 1 }} / {{ levels.length }}</span>
            <button
              class="btn-primary"
              @click="resetLevel"
              :style="{
                visibility: gameEnded ? 'hidden':'visible',
                fontSize: boardSizePx / 20 + 'px'
                }"
            >
              <font-awesome-icon icon="sync"></font-awesome-icon>
            </button>
          </div>
        </footer>
        <div class="fake-footer" v-if="isMobile"></div>
        <div class="link-footer" v-if="!isMobile">
          <a
            :href="'https://explorer2.harmony.one/#/address/' + globalData.address"
            target="_blank"
            class="link"
          >View Transactions</a>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Game from "./Game";
import Chip from "./Chip";
import StakeRow from "./StakeRow";
import TxHistoryLink from "./TxHistoryLink";
import RedeemPanel from "./RedeemPanel";
import { TweenLite } from "gsap/TweenMax";
import Vue from "vue";
import service from "../service";
import store from "../store";
import { levels } from "../level-generator";
import { setInterval, clearInterval } from "timers";

const InitialSeconds = 30;
function guid() {
  return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, function(c) {
    var r = (Math.random() * 16) | 0,
      v = c == "x" ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}

function getParameterByName(name) {
  var undefined;
  name = name.replace(/[\[]/, "\\[").replace(/[\]]/, "\\]");
  var regex = new RegExp("[\\?&]" + name + "=([^&#]*)"),
    results = regex.exec(location.search);
  return results === null
    ? undefined
    : decodeURIComponent(results[1].replace(/\+/g, " "));
}

export default {
  name: "PuzzlePage",
  components: {
    Game,
    Chip,
    StakeRow,
    TxHistoryLink,
    RedeemPanel
  },
  data() {
    return {
      globalData: store.data,
      levelIndex: 0,
      levels: [],
      boardSizePx: 0,
      size: 3,
      gameStarted: false,
      gameEnded: false,
      secondsLeft: InitialSeconds,
      timer: null,
      timeIncrease: "",
      balanceIncrease: "",
      isMobile: /iPhone|iPad|iPod|Android/i.test(navigator.userAgent),
      reward: 0,
      cancelEmail: false
    };
  },
  mounted() {
    let id = getParameterByName("cos");
    this.levels = levels();
    //Set board size follow height of screen when change screen size
    window.addEventListener(
      "resize",
      () => {
        this.boardSizePx = Math.min(
          this.$refs.gameContainer.clientWidth,
          window.innerHeight / 1.7
        );
        this.$forceUpdate;
      },
      false
    );
    // Set board size follow height of screen
    this.boardSizePx = Math.min(
      this.$refs.gameContainer.clientWidth,
      window.innerHeight / 1.7
    );
    service.register(id);
  },
  computed: {
    gameOverStyle() {
      return { fontSize: this.boardSizePx / 6 + "px" };
    },
    gameTutorialStyle() {
      return { fontSize: this.boardSizePx / 14 + "px" };
    },
    gameTutorialSmallStyle() {
      return { fontSize: this.boardSizePx / 16 + "px" };
    },
    infoItemStyle() {
      return { fontSize: this.boardSizePx / 18 + "px" };
    },
    levelTextStyle() {
      return { fontSize: this.boardSizePx / 18 + "px" };
    },
    stakeRowStyle() {
      return {
        width: this.boardSizePx + "px",
        fontSize: this.boardSizePx / 20 + "px"
      };
    },
    titleStyle() {
      return {
        fontSize: this.boardSizePx / 20 + "px"
      };
    },
    iconTokenStyle() {
      return {
        width: this.boardSizePx / 12 + "px",
        height: this.boardSizePx / 12 + "px"
      };
    },
    boardWrapperStyle() {
      return {
        width: this.boardSizePx + "px",
        height: this.boardSizePx + "px"
      };
    },
    level() {
      return this.levels[this.levelIndex];
    }
  },
  destroyed() {
    // Remove event change screen
    window.removeEventListener("resize", this.handleResize);
  },
  methods: {
    startGame() {
      this.gameStarted = true;
      this.gameEnded = false;
      this.cancelEmail = false;
      this.levelIndex = 0;
      this.reward = 0;
      this.levels = levels();
      this.secondsLeft = InitialSeconds;
      this.timer = setInterval(() => {
        this.secondsLeft--;
        if (this.secondsLeft <= 0) {
          this.endGame();
        }
      }, 1000);
    },
    resetLevel() {
      this.$refs[`game${this.levelIndex}`][0].reset();
    },
    onLevelComplete(moves) {

      if (this.levelIndex === this.levels.length - 1) {
        this.endGame();
        return;
      }
      service
        .completeLevel(this.globalData.account, this.levelIndex + 1, moves)
        .then(rewards => {
          this.levelIndex++;
          let timeChange = 15;
          this.secondsLeft += timeChange;
          this.timeIncrease = `+${timeChange}`;
          this.balanceIncrease = `+${rewards}`;
          this.reward += rewards;
          Vue.nextTick(() => {
            this.timeIncrease = "";
            this.balanceIncrease = "";
          });
        });
    },
    endGame() {
      stopBackgroundMusic();
      this.gameEnded = true;
      this.gameStarted = false;
      store.data.stake = 20;
      clearInterval(this.timer);
      this.timer = null;
    },
    restart() {
      this.gameEnded = false;
    },
    closeEmailPopup() {
      this.cancelEmail = true;
    }
  }
};
</script>
