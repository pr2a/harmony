<style scoped lang="less">
.score-container {
  margin: 0 auto 1em;
}

footer {
  margin: 1em auto 0;
  .btn-primary {
    font-size: 0.8em;
  }
}

.board-wrapper {
  position: relative;
  margin: 0 auto;
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
}

.main-container {
  .game-container {
    position: relative;
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
  .content {
    position: relative;
  }
}
.count-down {
  .seconds-left {
    &.game-over,
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
</style>

<template>
  <div id="app">
    <div class="main-container appearing">
      <div class="game-container" ref="gameContainer">
        <div class="score-container" :style="{ width: boardSizePx + 'px' }">
          <div class="logo"></div>
          <div class="flex-horizontal">
            <div class="count-down info-item">
              <div class="label">Time Left</div>
              <div class="content">
                <div
                  class="seconds-left"
                  :class="{ 'hurry-up': secondsLeft && secondsLeft <= 12, 'game-over': !secondsLeft }"
                >{{ secondsLeft }}</div>
                <transition>
                  <span v-if="timeIncrease!=''" class="number-increase">
                    {{
                    timeIncrease
                    }}
                  </span>
                </transition>
              </div>
            </div>
            <div class="balance info-item">
              <div class="label">Balance</div>
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
          </div>
        </div>

        <div class="board-wrapper" :style="boardWrapperStyle">
          <div v-if="gameEnded || !gameStarted">
            <div class="overlay game-over-message appearing">
              <div class="content">
                <p :style="gameOverStyle">{{ gameEnded ? 'Game over!' : 'Game Not Started' }}</p>
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
              :gameEnded="gameEnded"
              @completeLevel="onLevelComplete"
              v-if="i === levelIndex"
            ></Game>
          </transition>
        </div>
        <stake-row v-if="!gameStarted" @stake="startGame"></stake-row>
        <footer class="flex-vertical" :style="{ width: boardSizePx + 'px' }" v-if="gameStarted">
          <div class="flex-horizontal action-row">
            <span class="flex-grow">levels: {{ levelIndex + 1 }} / {{ levels.length }}</span>
            <button
              class="btn-primary"
              @click="resetLevel"
              :style="{ visibility: gameEnded ? 'hidden':'visible' }"
            >
              <font-awesome-icon icon="sync"></font-awesome-icon>
            </button>
          </div>
        </footer>
      </div>
    </div>
  </div>
</template>

<script>
import Game from "./Game";
import Chip from "./Chip";
import StakeRow from "./StakeRow";
import { TweenLite } from "gsap/TweenMax";
import Vue from "vue";
import service from "../service";
import store from "../store";
import { levels } from "../level-generator";
import { setInterval, clearInterval } from "timers";

const DefaultBoardSizePx = 420;
const InitialSeconds = 20;

export default {
  name: "PuzzlePage",
  components: {
    Game,
    Chip,
    StakeRow
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
      balanceIncrease: ""
    };
  },
  mounted() {
    this.levels = levels();
    // this.startGame();
    this.boardSizePx = Math.min(
      this.$refs.gameContainer.clientWidth,
      DefaultBoardSizePx
    );
  },
  computed: {
    gameOverStyle() {
      return { fontSize: this.boardSizePx / 6 + "px" };
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
  methods: {
    startGame() {
      this.gameStarted = true;
      this.gameEnded = false;
      this.levelIndex = 0;
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
      service.completeLevel(this.levelIndex, this.level, moves).then(() => {
        this.levelIndex++;
        let timeChange = 15;
        this.secondsLeft += timeChange;
        this.timeIncrease = `+${timeChange}`;
        let balanceChange = 4;
        // this.globalData.balance += balanceChange;
        this.balanceIncrease = `+${balanceChange}`;
        Vue.nextTick(() => {
          this.timeIncrease = "";
          this.balanceIncrease = "";
        });
      });
    },
    endGame() {
      this.gameEnded = true;
      this.gameStarted = false;
      clearInterval(this.timer);
      this.timer = null;
    },
    restart() {
      this.gameEnded = false;
    }
  }
};
</script>
