<template id="game2048" hidden>
  <div class="board" :tabindex="tabIndex" :style="boardStyle">
    <div ref="cells" v-for="(cl, i) in cells" class="cell" :key="i" :style="cellStyle">
      <Chip
        ref="chips"
        v-for="(ch, j) in cl.chips"
        :key="j"
        :animation-time-ms="animationTimeMs"
        :chip="ch"
        :size-px="cellSizePx"
      ></Chip>
    </div>
  </div>
</template>

<script>
import Chip from "./Chip";

function createGame2048(size) {
  size = size || 4;
  var size2 = size * size;
  var size1_2 = size / 2;
  var score = 0;

  var board = Array.apply(null, { length: size }).map(function() {
    return Array.apply(null, { length: size }).map(function() {
      return 0;
    });
  });
  console.log(board);

  function cellIsEmpty(c) {
    return board[c.y][c.x] == 0;
  }
  function cellsEqual(c1, c2) {
    return board[c1.y][c1.x] == board[c2.y][c2.x];
  }
  function canMoveChip(cf, ct) {
    return cellIsEmpty(ct) || cellsEqual(cf, ct);
  }

  function moveChip(cf, ct) {
    var tWasEmpty = cellIsEmpty(ct);
    var v = (board[ct.y][ct.x] += board[cf.y][cf.x]);
    board[cf.y][cf.x] = 0;
    return tWasEmpty ? 0 : v;
  }

  function findRandomEmptyPos() {
    var r = Math.floor(Math.random() * size2);
    var c = {};
    for (var i = size2; i > 0; i--) {
      c.y = Math.floor(r / size);
      c.x = r % size;
      if (cellIsEmpty(c)) return c;
      r++;
      if (r == size2) r = 0;
    }
    return null;
  }

  function rot0(c, x, y) {
    c.x = x;
    c.y = y;
  }
  function rot90(c, x, y) {
    c.x = y;
    c.y = x;
  }
  function rot180(c, x, y) {
    c.x = size - 1 - x;
    c.y = y;
  }
  function rot270(c, x, y) {
    c.x = y;
    c.y = size - 1 - x;
  }

  function move(rot) {
    var scoreInc = 0;
    var moves = [];
    var consolidations = [];
    var c = {};
    var tc = {};
    for (var y = 0; y < size; y++) {
      var s = size;
      for (var x = size - 2; x >= 0; x--) {
        rot(c, x, y);
        if (!cellIsEmpty(c)) {
          var tx = x;
          while (tx + 1 < s) {
            rot(tc, tx + 1, y);
            if (!cellIsEmpty(tc)) {
              if (cellsEqual(c, tc)) {
                tx++;
                s = tx;
              }
              break;
            }
            tx++;
          }

          if (x != tx) {
            rot(tc, tx, y);
            var v = moveChip(c, tc);
            moves.push({
              from: { x: c.x, y: c.y },
              to: { x: tc.x, y: tc.y }
            });
            if (v > 0) {
              consolidations.push({ x: tc.x, y: tc.y, value: v });
              scoreInc += v;
              score += v;
            }
          }
        }
      }
    }
    return {
      moves: moves,
      consolidations: consolidations,
      scoreInc: scoreInc
    };
  }

  return {
    size: size,
    board: board,
    score: function() {
      return score;
    },
    turn: function() {
      var chips = [];
      var p = findRandomEmptyPos();
      if (p != null) {
        var rnd = Math.random();
        var v = rnd > 0.8 ? 4 : 2;
        p.value = v;
        board[p.y][p.x] = v;
        chips.push(p);
      }
      return chips;
    },
    right: function() {
      return move(rot0);
    },
    down: function() {
      return move(rot90);
    },
    left: function() {
      return move(rot180);
    },
    up: function() {
      return move(rot270);
    },
    canMove: function() {
      for (
        var c = { y: 0 }, cr = { y: 0 }, cb = { y: 1 };
        c.y < size;
        c.y++, cr.y++, cb.y++
      )
        for (c.x = 0, cr.x = 1, cb.x = 0; c.x < size; c.x++, cr.x++, cb.x++) {
          if (
            cellIsEmpty(c) ||
            (cr.x < size && cellsEqual(c, cr)) ||
            (cb.y < size && cellsEqual(c, cb))
          )
            return true;
        }
    }
  };
}

function deffered(delayMs, func) {
  var executed = false;
  function execute() {
    if (!executed) {
      func();
      executed = true;
    }
  }
  function renew() {
    executed = false;
    setTimeout(execute, delayMs);
  }
  renew();
  return {
    finish: execute,
    renew: renew
  };
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
    attach: function(el) {
      el.addEventListener("touchstart", onStart, false);
      el.addEventListener("touchend", onEnd, false);
    },
    detach: function(el) {
      el.removeEventListener("touchstart", onStart);
      el.removeEventListener("touchend", onEnd);
    }
  };
}

var keyMap = {};
keyMap[37] = "left";
keyMap[38] = "up";
keyMap[39] = "right";
keyMap[40] = "down";

export default {
  name: "Game",
  components: {
    Chip
  },

  props: {
    size: { type: Number },
    sizeAimMap: { type: Array, required: true },
    listenOwnKeyEventsOnly: { type: Boolean, default: false },
    tabIndex: { type: Number, default: 1 },
    boardSizePx: { type: Number, default: 0 },
    animationTimeMs: { type: Number, default: 150 },
    started: { type: Boolean, default: false }
  },
  data: function() {
    var aim = this.sizeAimMap[this.size];
    return {
      cells: this.createCells(),
      boardSizeAutoPx: 0,
      aim: aim
    };
  },
  mounted: function() {
    this.boardSizeAutoPx =
      this.boardSizePx > 0
        ? this.boardSizePx
        : this.$el.getBoundingClientRect().width;
  },
  watch: {
    size: function() {
      this.cells = this.createCells();
      this.aim = this.sizeAimMap[this.size];
      this.$emit("aim-changed", this.aim);
    },
    started: function(nv, ov) {
      if (nv) {
        this.startGame();
      } else {
        this.endGame();
      }
    }
  },
  computed: {
    boardStyle: function() {
      return {
        width: this.boardSizePx > 0 ? this.boardSizePx + "px" : "100%",
        height: this.boardSizePx > 0 ? this.boardSizePx + "px" : "100%",
        borderRadius: 7 / this.size + "%"
      };
    },
    cellStyle: function() {
      return {
        width: this.cellSizePct + "%",
        height: this.cellSizePct + "%",
        marginLeft: this.cellMarginPct + "%",
        marginTop: this.cellMarginPct + "%"
      };
    },
    cellSizePct: function() {
      return 8 * this.cellMarginPct;
    },
    cellMarginPct: function() {
      return 100 / (9 * this.size + 1);
    },
    cellSizePx: function() {
      return (
        (this.cellSizePct / 100) *
        (this.boardSizePx > 0 ? this.boardSizePx : this.boardSizeAutoPx)
      );
    }
  },
  methods: {
    startGame: function() {
      this.emptyCells();
      var game = createGame2048(this.size);
      for (var i = Math.max(2, this.size - 2); i > 0; i--) {
        var chips = game.turn();
        console.log("start game, chips", chips);
        this.addChips(chips);
      }
      var doGameMove = this.createGameMove(game);
      this.runKeyboardControl(doGameMove);
      this.runTouchControl(doGameMove);
      this.$emit("started", this);
    },

    runKeyboardControl: function(doGameMove) {
      var listenKeysOn = this.listenOwnKeyEventsOnly ? this.$el : document;
      var h = function(e) {
        var m = keyMap[e.keyCode];
        if (m == null) return;
        e.preventDefault();
        doGameMove(m);
      };
      listenKeysOn.addEventListener("keydown", h);
      this.$once("ended", function() {
        listenKeysOn.removeEventListener("keydown", h);
      });
    },

    runTouchControl: function(doGameMove) {
      var sw = createSwipeListener(function(m) {
        doGameMove(m);
      });
      var el = this.$el;
      sw.attach(el);
      this.$once("ended", function() {
        sw.detach(el);
      });
    },

    createGameMove: function(game) {
      var self = this;
      var boardChanges = { consolidations: [] };
      var newChips = [];
      var consolidateAndAddChipsDeffered = deffered(
        self.animationTimeMs,
        function() {
          self.consolidateChips(boardChanges.consolidations);
          self.addChips(newChips);
        }
      );

      return function(m) {
        consolidateAndAddChipsDeffered.finish();

        boardChanges = game[m]();
        newChips.length = 0;
        if (boardChanges.moves.length > 0) {
          for (var i = Math.max(1, self.size - 3); i > 0; i--) {
            var chips = game.turn();
            chips.push.apply(newChips, chips);
          }
          if (boardChanges.scoreInc > 0) {
            self.$emit("score", {
              score: game.score(),
              scoreInc: boardChanges.scoreInc
            });
            for (var i = 0; i < boardChanges.consolidations.length; i++) {
              if (boardChanges.consolidations[i].value == self.aim) {
                self.$emit("aim-reached");
                break;
              }
            }
          }
        }

        self.moveChips(boardChanges.moves);
        consolidateAndAddChipsDeffered.renew();
        if (!game.canMove()) {
          setTimeout(function() {
            self.endGame();
          }, self.animationTimeMs);
        }
      };
    },

    endGame: function() {
      this.$emit("ended", this);
    },

    consolidateChips: function(consolidations) {
      var self = this;
      consolidations.forEach(function(c) {
        var cell = self.getCell(c);
        var chips = cell.chips;
        chips.splice(0, chips.length - 1);
        chips[0].value = c.value;
      });
    },
    moveChips: function(moves) {
      for (var i = 0; i < moves.length; i++)
        this.moveChip(moves[i].from, moves[i].to);
    },
    moveChip: function(from, to) {
      var fcell = this.getCell(from);
      var fcellEl = this.getCellEl(from);
      var tcell = this.getCell(to);
      var tcellEl = this.getCellEl(to);
      var chip = fcell.chips.splice(0, 1)[0];
      var fboundRect = fcellEl.getBoundingClientRect();
      var tboundRect = tcellEl.getBoundingClientRect();
      chip.prevRelPos = {
        left: fboundRect.left - tboundRect.left,
        top: fboundRect.top - tboundRect.top
      };
      tcell.chips.push(chip);
    },
    addChips: function(chips) {
      console.log("minh1", chips);
      chips.forEach(this.addChip);
    },
    addChip: function(c) {
      console.log("minh2", c);
      this.cells[this.getCellIndex(c)].chips.push({ value: c.value });
      console.log("minh3", this.cells[this.getCellIndex(c)].chips);
    },
    getCellIndex: function(c) {
      return c.y * this.size + c.x;
    },
    getCell: function(c) {
      return this.cells[this.getCellIndex(c)];
    },
    getCellEl: function(c) {
      return this.$refs.cells[this.getCellIndex(c)];
    },
    createCells: function() {
      return Array.apply(null, { length: this.size * this.size }).map(
        function() {
          return { chips: [] };
        }
      );
    },
    emptyCells: function() {
      this.cells.forEach(function(c) {
        c.chips.splice(0);
      });
    }
  }
};
</script>
