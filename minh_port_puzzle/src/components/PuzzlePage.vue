<style scoped lang="less">
.score-container {
  margin-bottom: 1em;
}

.info-item {
  background-color: #fff;
  border-radius: 0.5em;
  border: 0.15em solid #979797;
  padding: 1em;
  border-radius: 0.25em;
  padding: 0.5em;
  width: 6em;
  margin-left: 1em;
  text-align: center;
  .label {
    font-size: 0.5em;
    margin-bottom: 0.5em;
  }
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

.game-over {
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
            <div class="content">{{ secondsLeft }}</div>
          </div>
          <div class="reward info-item">
            <div class="label">Reward</div>
            <div class="content">{{ reward }}</div>
          </div>
        </div>
      </div>
      <div class="game-container" :style="gameContainerStyle">
        <div v-if="gameEnded">
          <div class="overlay half-white appearing07"></div>
          <div class="overlay game-over appearing">
            <div class="content">
              <p :style="gameOverStyle">Game over!</p>
              <div>
                <button class="btn-primary" @click="startGame">Restart</button>
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
import { levels } from "../level-generator";
import { setInterval, clearInterval } from "timers";
import TxHistoryPanel from "./TxHistoryPanel";

var defBoardSizePx = 420;
var defSize = 3;

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
      levelIndex: 0,
      levels: [],
      boardSizePx: defBoardSizePx,
      size: defSize,
      gameEnded: false,
      secondsLeft: IntialSeconds,
      reward: 0,
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
      service.stakeToken().then(() => {
        this.timer = setInterval(() => {
          this.secondsLeft--;
          if (this.secondsLeft <= 0) {
            this.endGame();
          }
        }, 1000);
      });
    },
    resetLevel() {
      this.$refs[`game${this.levelIndex}`][0].reset();
    },
    onLevelComplete() {
      if (this.levelIndex === this.levels.length - 1) {
        this.endGame();
        return;
      }
      service.completeLevel(this.levelIndex).then(() => {
        this.levelIndex++;
        this.secondsLeft += 15;
        this.reward += 5;
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
