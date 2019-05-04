<style scoped lang="less">
.board {
  display: flex;
  flex-wrap: wrap;
  align-content: flex-start;
  align-items: space-around;
  background-color: white;
  outline: none;
  position: absolute;
  border-radius: 0.5em;
  margin: 0 auto;

  .cell {
    background-color: #ada49f;
    position: relative;
    border-radius: 7%;
    &.selected {
      box-shadow: 0 0 0 0.4em rgba(0, 0, 0, 0.8);
    }
    .chip {
      position: absolute;
      width: 100%;
      height: 100%;
      display: flex;
      justify-content: center;
      align-items: center;
      overflow: hidden;
      text-align: justify;
      font-weight: bold;
      background-color: honeydew;
      z-index: 1;
      border-radius: 7%;
      /* Hack to improve transition performance on mobile devices. It enables GPU rendering. */
      transform: translateZ(0);
      -webkit-transform: translateZ(0);
    }
  }
}

.demo-arrow-1 {
  position: absolute;
  z-index: 2;
  display: inline-block;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  opacity: 0.5;
  background-size: cover;
  background-image: url(../assets/pointer-dark.png);
  transform: translateX(5%);
}

.demo-arrow-2 {
  position: absolute;
  z-index: 2;
  display: inline-block;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  opacity: 0.5;
  background-size: cover;
  background-image: url(../assets/arrow-loop-dark-further.png);
  transform: translateX(5%);
}
</style>

<template>
  <div class="board" :tabindex="tabIndex" :style="boardStyle">
    <div v-if="gameLevel === 1 && gameStarted" class="demo-arrow-1"></div>
    <div
      ref="cells"
      v-for="(value, i) in cells"
      class="cell"
      :class="{ selected: index === i }"
      :key="i"
      :style="cellStyle"
    >
      <Chip ref="chips" :animation-time-ms="animationTimeMs" :value="value" :size-px="cellSizePx"></Chip>
    </div>
  </div>
</template>

<script>
import Chip from "./Chip";
import Vue from "vue";
import { playMoveSound, playBeginSound, playEndSound } from "../lib/sound";
import { constants } from "fs";

function clamp(v, min, max) {
  return Math.max(min, Math.min(max, v));
}

function createSwipeListener(onSwipe) {
  var sens = 5;
  var st;

  function onStart(e) {
    st = e.touches[0];
    e.preventDefault();
  }

  function onEnd(e) {
    var et = e.changedTouches[0];
    var x = st.clientX - et.clientX;
    var y = st.clientY - et.clientY;
    var mx = Math.abs(x);
    var my = Math.abs(y);
    if (mx < sens && my < sens) return;

    var d = mx > my ? (x > 0 ? "L" : "R") : y > 0 ? "U" : "D";
    onSwipe(d);
  }

  return {
    attach(el) {
      el.addEventListener("touchstart", onStart, false);
      el.addEventListener("touchend", onEnd, false);
    },
    detach(el) {
      el.removeEventListener("touchstart", onStart);
      el.removeEventListener("touchend", onEnd);
    }
  };
}

const actions = {
  L: { x: 0, y: -1 },
  U: { x: -1, y: 0 },
  R: { x: 0, y: 1 },
  D: { x: 1, y: 0 }
};
var keyMap = {};
keyMap[37] = "L";
keyMap[38] = "U";
keyMap[39] = "R";
keyMap[40] = "D";
keyMap[65] = "L";
keyMap[87] = "U";
keyMap[68] = "R";
keyMap[83] = "D";

export default {
  name: "Game",
  components: {
    Chip
  },

  props: {
    game: { contents: Array, initialSelected: { x: Number, y: Number } },
    gameLevel: Number,
    listenOwnKeyEventsOnly: { type: Boolean, default: false },
    tabIndex: { type: Number, default: 1 },
    boardSizePx: { type: Number, default: 0 },
    animationTimeMs: { type: Number, default: 150 },
    gameEnded: Boolean,
    gameStarted: Boolean
  },
  data() {
    return {
      cells: this.game.contents.slice(0),
      position: Object.assign({}, this.game.initialSelected),
      boardSizeAutoPx: 0,
      size: 3,
      moves: ""
    };
  },
  mounted() {
    this.boardSizeAutoPx =
      this.boardSizePx > 0
        ? this.boardSizePx
        : this.$el.getBoundingClientRect().width;
    this.startGame();
  },
  computed: {
    index() {
      return this.position.x * 3 + this.position.y;
    },
    boardStyle() {
      return {
        width: this.boardSizePx > 0 ? this.boardSizePx + "px" : "100%",
        height: this.boardSizePx > 0 ? this.boardSizePx + "px" : "100%",
        borderRadius: 7 / this.size + "%"
      };
    },
    cellStyle() {
      return {
        width: this.cellSizePct + "%",
        height: this.cellSizePct + "%",
        marginLeft: this.cellMarginPct + "%",
        marginTop: this.cellMarginPct + "%"
      };
    },
    cellSizePct() {
      return 8 * this.cellMarginPct;
    },
    cellMarginPct() {
      return 100 / (9 * this.size + 1);
    },
    cellSizePx() {
      return (
        (this.cellSizePct / 100) *
        (this.boardSizePx > 0 ? this.boardSizePx : this.boardSizeAutoPx)
      );
    }
  },
  watch: {
    gameEnded(val) {
      if (val) {
        this.$emit("ended");
      }
    }
  },
  methods: {
    startGame() {
      // Add begin sound.
      playBeginSound();
      this.runKeyboardControl(this.move);
      this.runTouchControl(this.move);
    },
    runKeyboardControl(move) {
      var listenKeysOn = this.listenOwnKeyEventsOnly ? this.$el : document;
      var h = e => {
        if (!this.gameStarted) return;
        var m = keyMap[e.keyCode];
        if (m == null) return;
        e.preventDefault();
        // Add sound before any move.
        playMoveSound();
        move(m);
      };
      listenKeysOn.addEventListener("keydown", h);
      // TODO: on game end, remove listeners.
      this.$once("completeLevel", function() {
        listenKeysOn.removeEventListener("keydown", h);
      });
    },

    runTouchControl(move) {
      var sw = createSwipeListener(m => {
        if (!this.gameStarted) return;
        // Add sound before any move.
        playMoveSound();
        move(m);
      });
      var el = this.$el;
      sw.attach(el);
      this.$once("completeLevel", function() {
        sw.detach(el);
      });
    },
    finishLevel() {
      this.$emit("completeLevel", this.moves);
    },
    move(dir) {
      this.moves += dir;
      let diff = actions[dir];
      let x = clamp(this.position.x + diff.x, 0, 2);
      let y = clamp(this.position.y + diff.y, 0, 2);
      if (x === this.position.x && y === this.position.y) return;
      this.position.x = x;
      this.position.y = y;
      this.cells[this.index]++;
      if (this.isLevelPassed()) {
        this.finishLevel();
      }
    },
    isLevelPassed() {
      let v = this.cells[0];
      return this.cells.findIndex(x => x !== v) === -1;
    },
    reset() {
      this.cells = this.game.contents.slice(0);
      this.position = Object.assign(
        {},
        this.position,
        this.game.initialSelected
      );
    }
  }
};
</script>
