<style scoped lang="less">
.score-container {
  margin-bottom: 1em;
}

footer {
  margin-top: 1em;
  .btn-primary {
    font-size: 0.8em;
  }
}

.game-wrapper {
  position: relative;
  .game {
    position: absolute;
  }
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

.tx-history-panel {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1000;
  background-color: rgba(252, 247, 235, 0.95);
}
.view-tx-btn {
  font-size: 0.8em;
}
.action-row + .action-row {
  margin-top: 1em;
}
.count-down .content {
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
</style>

<template>
  <div id="app" style="visibility:hidden">
    <tx-history-panel v-if="isTxPanelOpen" class="tx-history-panel" @close="isTxPanelOpen = false"></tx-history-panel>
    <div class="main-container appearing" :style="mainContainerStyle">
      <div class="score-container">
        <div class="logo"></div>
        <div class="flex-horizontal">
          <div class="count-down info-item">
            <div class="label">Time Left</div>
            <div
              class="content"
              :class="{ 'hurry-up': secondsLeft && secondsLeft <= 12, 'game-over': !secondsLeft }"
            >{{ secondsLeft }}</div>
          </div>
          <div class="balance info-item">
            <div class="label">Balance</div>
            <div class="content">{{ globalData.balance }}</div>
          </div>
        </div>
      </div>
      <div class="game-container" :style="gameContainerStyle">
        <div v-if="gameEnded">
          <div class="overlay half-white appearing07"></div>
          <div class="overlay game-over-message appearing">
            <div class="content">
              <p :style="gameOverStyle">Game over!</p>
              <div>
                <button class="btn-primary" @click="$emit('restart')">Restart</button>
              </div>
            </div>
          </div>
        </div>

        <div class="game-wrapper" :style="boardStyle">
          <transition name="fade" v-for="(level, i) in levels" :key="i">
            <Game
              :ref="'game' + i"
              class="game"
              :listen-own-key-events-only="false"
              :tab-index="1"
              :board-size-px="boardSizePx"
              :game="level"
              :gameEnded="gameEnded"
              @completeLevel="onLevelComplete"
              v-if="i === levelIndex"
            ></Game>
          </transition>
        </div>

        <footer class="flex-vertical">
          <div class="flex-horizontal action-row">
            <span class="flex-grow">levels: {{ levelIndex + 1 }} / {{ levels.length }}</span>
            <button
              class="btn-primary"
              @click="resetLevel"
              :style="{ visibility: gameEnded ? 'hidden':'visible' }"
            >Reset Level</button>
          </div>
          <div class="flex-horizontal action-row">
            <div class="flex-grow"></div>
            <button class="btn-primary view-tx-btn" @click="viewTxHistory">View Transactions</button>
          </div>
        </footer>
      </div>
    </div>
  </div>
</template>

<script>
import Game from "./Game";
import Chip from "./Chip";
import { TweenLite } from "gsap/TweenMax";
import Vue from "vue";
import service from "../service";
import store from "../store";
import { levels } from "../level-generator";
import { setInterval, clearInterval } from "timers";
import TxHistoryPanel from "./TxHistoryPanel";

var defBoardSizePx = 420;

const IntialSeconds = 30;
export default {
  name: "PuzzlePage",
  components: {
    Game,
    Chip,
    TxHistoryPanel
  },
  data() {
    return {
      globalData: store.data,
      levelIndex: 0,
      levels: [],
      boardSizePx: defBoardSizePx,
      size: 3,
      gameEnded: false,
      secondsLeft: IntialSeconds,
      timer: null,
      isTxPanelOpen: false
    };
  },
  created() {
    this.loadState();
  },
  mounted() {
    var self = this;
    this.startGame();
    requestAnimationFrame(function() {
      self.fitBoardSizePx();
      requestAnimationFrame(function() {
        self.$el.style.visibility = "visible";
      });
    });
  },
  computed: {
    boardStyle() {
      return {
        width: this.boardSizePx > 0 ? this.boardSizePx + "px" : "100%",
        height: this.boardSizePx > 0 ? this.boardSizePx + "px" : "100%",
        borderRadius: 7 / this.size + "%"
      };
    },
    gameOverStyle() {
      return { fontSize: this.boardSizePx / 6 + "px" };
    },
    gameContainerStyle() {
      return {
        width: this.boardSizePx + "px",
        height: this.boardSizePx + "px"
      };
    },
    mainContainerStyle() {
      return {
        width: this.boardSizePx + "px"
      };
    },
    level() {
      return this.levels[this.levelIndex];
    }
  },
  methods: {
    fitBoardSizePx() {
      if (window.innerWidth < defBoardSizePx * 1.04) {
        this.boardSizePx = window.innerWidth * 0.96;
      } else {
        this.boardSizePx = defBoardSizePx;
      }
    },
    loadState() {
      try {
        var s = document.cookie;
        if (s) {
          var state = JSON.parse(s);
          if (state) {
            if (state.awards) this.awards = state.awards;
            if (state.bestScore) this.bestScore = state.bestScore;
          }
        }
      } catch (e) {}
    },
    persistState() {
      try {
        var state = {
          bestScore: this.bestScore,
          awards: this.awards
        };
        document.cookie = JSON.stringify(state);
      } catch (e) {}
    },
    startGame() {
      this.gameEnded = false;
      this.levelIndex = 0;
      this.levels = levels();
      this.secondsLeft = IntialSeconds;
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
      service.completeLevel(this.levelIndex, moves).then(() => {
        this.levelIndex++;
        this.secondsLeft += 15;
        this.balance += 5;
        this.persistState();
      });
    },
    endGame() {
      this.gameEnded = true;
      clearInterval(this.timer);
      this.timer = null;
    },
    viewTxHistory() {
      this.isTxPanelOpen = true;
    }
  }
};
</script>
