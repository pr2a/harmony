<template>
  <div id="app" style="visibility:hidden">
    <div class="main-container appearing" :style="mainContainerStyle">
      <div class="score-container" :style="scoreContainerStyle">
        <div
          ref="gameAim"
          class="game-aim"
          v-bind:class="{'game-aim-reached':gameAimReached}"
          :style="gameAimStyle"
        >{{ gameAim }}</div>
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
      <div class="game-controls" :style="gameControlsStyle">
        <div class="size-control" v-if="!gameStarted">
          Size:
          <div class="size-control-item" v-for="s in sizes" :key="s">
            <input type="radio" :id="'size-radio'+s" :value="s" v-model.number="size">
            <label :for="'size-radio'+s">{{ s }}</label>
          </div>&nbsp;
        </div>
        <button
          v-if="!gameStarted"
          @click="startGame()"
          class="button"
          :style="buttonStyle"
          key="start"
        >New Game</button>
        <button v-else @click="gameStarted=false" class="button" :style="buttonStyle" key="end">End</button>
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
  name: "Puzzle",
  components: {
    Game,
    Chip
  },
  data: function() {
    var sizeAimMap = [];
    sizeAimMap[3] = 256;
    // sizeAimMap[4] = 2048;
    // sizeAimMap[5] = 4096;
    // sizeAimMap[6] = 8192;

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
  created: function() {
    this.loadState();
  },
  mounted: function() {
    var self = this;
    requestAnimationFrame(function() {
      self.fitBoardSizePx();
      requestAnimationFrame(function() {
        self.$el.style.visibility = "visible";
      });
    });
  },
  computed: {
    gameOverStyle: function() {
      return { fontSize: this.boardSizePx / 6 + "px" };
    },
    gameContainerStyle: function() {
      return {
        width: this.boardSizePx + "px",
        height: this.boardSizePx + "px"
      };
    },
    mainContainerStyle: function() {
      return {
        width: this.boardSizePx + "px"
      };
    },
    gameControlsStyle: function() {
      return {
        height: this.boardSizePx * 0.2 + "px"
      };
    },
    scoreContainerStyle: function() {
      return {
        height: this.boardSizePx * 0.2 + "px"
      };
    },
    gameAimStyle: function() {
      var bsh = this.boardSizePx / 50 + "px ";
      return {
        boxShadow: "0 " + bsh + bsh + "black",
        fontSize: this.boardSizePx / 110 + "em"
      };
    },
    buttonStyle: function() {
      return {
        fontSize: this.boardSizePx / 450 + "em"
      };
    },
    scoreStyle: function() {
      return {
        fontSize: this.boardSizePx / 280 + "em"
      };
    },
    gameAwardsContainerStyle: function() {
      return {
        height: this.boardSizePx * 0.08 + "px"
      };
    },
    gameAwardStyle: function() {
      return {
        width: this.boardSizePx / 5 + "px",
        fontSize: this.boardSizePx / 350 + "em"
      };
    },
    gameAwardLikeStyle: function() {
      return {
        height: this.boardSizePx / 21 + "px"
      };
    },
    allAwardsObtained: function() {
      for (var i in this.awards) if (!this.awards[i].obtained) return false;
      return true;
    }
  },
  watch: {
    size: function() {
      this.gameEnded = false;
    }
  },
  methods: {
    fitBoardSizePx: function() {
      if (window.innerWidth < defBoardSizePx * 1.04) {
        this.boardSizePx = window.innerWidth * 0.96;
      } else {
        this.boardSizePx = defBoardSizePx;
      }
    },
    loadState: function() {
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
    persistState: function() {
      try {
        var state = {
          bestScore: this.bestScore,
          awards: this.awards
        };
        document.cookie = JSON.stringify(state);
      } catch (e) {}
    },
    startGame: function() {
      this.gameStarted = true;
      this.score = 0;
    },
    onGameStarted: function() {
      this.gameStarted = true;
      this.gameEnded = false;
    },
    onGameEnded: function() {
      this.gameStarted = false;
      this.gameEnded = true;
      this.gameAimReached = false;
      this.persistState();
    },
    onGameScore: function(args) {
      var s = { score: this.score };
      var self = this;
      TweenLite.to(s, 0.3, {
        score: args.score,
        ease: Power0.easeNone,
        onUpdate: function() {
          self.score = Math.floor(s.score);
        }
      });

      if (args.score > this.bestScore[this.size]) {
        var bs = { score: this.bestScore[this.size] };
        TweenLite.to(bs, 0.3, {
          score: args.score,
          ease: Power0.easeNone,
          onUpdate: function() {
            Vue.set(self.bestScore, self.size, Math.floor(bs.score));
          }
        });
      }

      this.scoreInc = args.scoreInc + "+";
      Vue.nextTick(function() {
        self.scoreInc = "";
      });
    },
    onGameAimChanged: function(aim) {
      this.gameAim = aim;
    },
    onGameAimReached: function() {
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
