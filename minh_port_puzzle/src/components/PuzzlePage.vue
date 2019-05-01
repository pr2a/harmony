<style scoped lang="less">
.score-container {
  margin-bottom: 1em;
}

.flex-horizontal {
  display: flex;
  flex-direction: row;
  align-items: center;
}

.flex-grow {
  flex: 1;
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
</style>

<template>
  <div id="app" style="visibility:hidden">
    <div class="main-container appearing" :style="mainContainerStyle">
      <div class="score-container">
        <div class="logo"></div>
      </div>
      <div class="game-container" :style="gameContainerStyle">
        <div v-if="gameEnded">
          <div class="overlay half-white appearing07"></div>
          <div class="overlay game-over appearing" :style="gameOverStyle">
            <p>Game over!</p>
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
              :started="gameStarted"
              :game="level"
              @ended="onGameEnded"
              @score="onGameScore"
              v-if="i === levelIndex"
            ></Game>
          </transition>
        </div>

        <footer class="flex-horizontal">
          <span class="flex-grow">levels: {{ levelIndex + 1 }} / {{ levels.length }}</span>
          <button class="btn-primary pull-right" @click="reset">Reset</button>
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
import { levels } from "../level-generator";

var defBoardSizePx = 420;
var defSize = 3;

export default {
  name: "PuzzlePage",
  components: {
    Game,
    Chip
  },
  data() {
    return {
      levelIndex: 0,
      levels: levels(),
      boardSizePx: defBoardSizePx,
      size: defSize,
      gameStarted: false,
      gameEnded: false,
      score: 0,
      scoreInc: ""
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
      this.gameStarted = true;
      this.score = 0;
    },
    reset() {
      console.log(this.$refs[`game${this.levelIndex}`]);
      this.$refs[`game${this.levelIndex}`][0].reset();
    },
    onGameEnded() {
      this.gameStarted = false;
      if (this.levelIndex === this.levels.length - 1) {
        this.gameEnded = true;
        return;
      }
      this.levelIndex++;
      this.persistState();
    },
    onGameScore(args) {
      var s = { score: this.score };
      var self = this;
      TweenLite.to(s, 0.3, {
        score: args.score,
        ease: Power0.easeNone,
        onUpdate() {
          self.score = Math.floor(s.score);
        }
      });

      if (args.score > this.bestScore[this.size]) {
        var bs = { score: this.bestScore[this.size] };
        TweenLite.to(bs, 0.3, {
          score: args.score,
          ease: Power0.easeNone,
          onUpdate() {
            Vue.set(self.bestScore, self.size, Math.floor(bs.score));
          }
        });
      }

      this.scoreInc = args.scoreInc + "+";
      Vue.nextTick(function() {
        self.scoreInc = "";
      });
    }
  }
};
</script>
