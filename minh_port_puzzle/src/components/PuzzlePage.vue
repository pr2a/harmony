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
</style>

<template>
  <div id="app" style="visibility:hidden">
    <div class="main-container appearing" :style="mainContainerStyle">
      <div class="score-container" :style="scoreContainerStyle">
        <!-- <div
          ref="gameAim"
          class="game-aim"
          v-bind:class="{'game-aim-reached':gameAimReached}"
          :style="gameAimStyle"
        >{{ gameAim }}</div>-->
        <div class="logo"></div>
        <div class="scores" :style="scoreStyle">
          <div class="score">
            <div class="label">Score</div>
            <div>
              {{ score }}
              <transition>
                <span v-if="scoreInc!=''" class="score-inc">
                  {{
                  scoreInc
                  }}
                </span>
              </transition>
            </div>
          </div>&nbsp;
          <div class="score">
            <div class="label">Best</div>
            <div>{{ bestScore[size] }}</div>
          </div>
        </div>
      </div>
      <div class="game-container" :style="gameContainerStyle">
        <div v-if="gameEnded">
          <div class="overlay half-white appearing07"></div>
          <div class="overlay game-over appearing" :style="gameOverStyle">
            <p>Game over!</p>
          </div>
        </div>
        <Game
          ref="game"
          :size="size"
          :size-aim-map="sizeAimMap"
          :listen-own-key-events-only="false"
          :tab-index="1"
          :board-size-px="boardSizePx"
          :started="gameStarted"
          @started="onGameStarted"
          @ended="onGameEnded"
          @score="onGameScore"
          @aim-changed="onGameAimChanged"
          @aim-reached="onGameAimReached"
        ></Game>
        <footer class="flex-horizontal">
          <span class="flex-grow">levels: {{ level }} / 100</span>
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

var defBoardSizePx = 420;
var defSize = 3;

export default {
  name: "PuzzlePage",
  components: {
    Game,
    Chip
  },
  data() {
    var sizeAimMap = [];
    sizeAimMap[3] = 256;

    var awards = {};
    var bestScore = {};
    var sizes = [];
    var i = 0;
    for (var s in sizeAimMap) {
      var a = sizeAimMap[s];
      bestScore[s] = 0;
      awards[a] = { aim: a, obtained: false };
      sizes[i++] = s;
    }

    return {
      level: 1,
      boardSizePx: defBoardSizePx,
      size: defSize,
      sizes: sizes,
      sizeAimMap: sizeAimMap,
      gameStarted: false,
      gameEnded: false,
      gameAim: sizeAimMap[defSize],
      gameAimReached: false,
      score: 0,
      scoreInc: "",
      bestScore: bestScore,
      awards: awards
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
    gameControlsStyle() {
      return {
        height: this.boardSizePx * 0.2 + "px"
      };
    },
    scoreContainerStyle() {
      return {
        height: this.boardSizePx * 0.2 + "px"
      };
    },
    gameAimStyle() {
      var bsh = this.boardSizePx / 50 + "px ";
      return {
        boxShadow: "0 " + bsh + bsh + "black",
        fontSize: this.boardSizePx / 110 + "em"
      };
    },
    scoreStyle() {
      return {
        fontSize: this.boardSizePx / 280 + "em"
      };
    },
    gameAwardsContainerStyle() {
      return {
        height: this.boardSizePx * 0.08 + "px"
      };
    },
    gameAwardStyle() {
      return {
        width: this.boardSizePx / 5 + "px",
        fontSize: this.boardSizePx / 350 + "em"
      };
    },
    gameAwardLikeStyle() {
      return {
        height: this.boardSizePx / 21 + "px"
      };
    },
    allAwardsObtained() {
      for (var i in this.awards) if (!this.awards[i].obtained) return false;
      return true;
    }
  },
  watch: {
    size() {
      this.gameEnded = false;
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
      throw "not implemented";
    },
    onGameStarted() {
      this.gameStarted = true;
      this.gameEnded = false;
    },
    onGameEnded() {
      this.gameStarted = false;
      this.gameEnded = true;
      this.gameAimReached = false;
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
    },
    onGameAimChanged(aim) {
      this.gameAim = aim;
    },
    onGameAimReached() {
      this.gameAimReached = true;
      this.awards[this.gameAim].obtained = true;
      this.persistState();

      var awardEl = this.getAwardEl(this.gameAim);
      var gameAimEl = this.$refs.gameAim;
      var p1 = gameAimEl.getBoundingClientRect();
      var p2 = awardEl.getBoundingClientRect();
      var ws = p1.width / p2.width;
      var hs = p1.height / p2.height;
      var x = p1.left - p2.left + p1.width / 4;
      var y = p1.top - p2.top + p1.height / 2;

      var s = awardEl.style;
      s["-webkit-transform"] = s.transform =
        "translate(" + x + "px," + y + "px) scale(" + ws + "," + hs + ")";
      s["-webkit-transition"] = s.transition = "";
      s.zIndex = 100;
      requestAnimationFrame(function() {
        s["-webkit-transition"] = s.transition = "all 2s";
        s["-webkit-transform"] = s.transform = "";
      });
    }
  }
};
</script>
