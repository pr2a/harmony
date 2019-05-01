<style scoped lang="less">
.board {
  border-radius: 0.5em;

  .cell.selected {
    box-shadow: 0 0 0 0.4em rgba(255, 255, 255, 0.4);
  }
}
</style>

<template>
  <div class="board" :tabindex="tabIndex" :style="boardStyle">
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

    var d = mx > my ? (x > 0 ? "left" : "right") : y > 0 ? "up" : "down";
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

var keyMap = {};
keyMap[37] = { dir: "left", diff: { x: 0, y: -1 } };
keyMap[38] = { dir: "up", diff: { x: -1, y: 0 } };
keyMap[39] = { dir: "right", diff: { x: 0, y: 1 } };
keyMap[40] = { dir: "down", diff: { x: 1, y: 0 } };

export default {
  name: "Game",
  components: {
    Chip
  },

  props: {
    game: { contents: Array, initialSelected: { x: Number, y: Number } },
    listenOwnKeyEventsOnly: { type: Boolean, default: false },
    tabIndex: { type: Number, default: 1 },
    boardSizePx: { type: Number, default: 0 },
    animationTimeMs: { type: Number, default: 150 }
  },
  data() {
    return {
      cells: this.game.contents,
      position: this.game.initialSelected,
      boardSizeAutoPx: 0,
      size: 3
    };
  },
  mounted() {
    console.log("xxx", this.game);
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
  methods: {
    startGame() {
      this.$emit("started", this);
      this.runKeyboardControl(this.move);
    },

    runKeyboardControl(move) {
      var listenKeysOn = this.listenOwnKeyEventsOnly ? this.$el : document;
      var h = function(e) {
        var m = keyMap[e.keyCode];
        if (m == null) return;
        e.preventDefault();
        move(m);
      };
      listenKeysOn.addEventListener("keydown", h);
      this.$once("ended", function() {
        listenKeysOn.removeEventListener("keydown", h);
      });
    },

    runTouchControl(doGameMove) {
      var sw = createSwipeListener(function(m) {
        doGameMove(m);
      });
      var el = this.$el;
      sw.attach(el);
      this.$once("ended", function() {
        sw.detach(el);
      });
    },
    endGame() {
      console.log("eee");
      this.$emit("ended", this);
    },
    move(e) {
      console.log("xx", this.position);
      let x = clamp(this.position.x + e.diff.x, 0, 2);
      let y = clamp(this.position.y + e.diff.y, 0, 2);
      if (x === this.position.x && y === this.position.y) return;
      this.position.x = x;
      this.position.y = y;
      this.cells[this.index]++;
      if (this.isGameEnded()) {
        this.endGame();
      }
    },
    isGameEnded() {
      let v = this.cells[0];
      return this.cells.findIndex(x => x !== v) === -1;
    }
  }
};
</script>
