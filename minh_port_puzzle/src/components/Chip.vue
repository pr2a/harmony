<template>
  <transition :css="false" @enter="enter">
    <div class="chip" :style="style">{{ chip.value }}</div>
  </transition>
</template>

<script>
var fontSizeCoefs = [1, 1, 0.8, 0.65, 0.5, 0.4, 0.35, 0.32];
var backColors = [];
backColors[2] = "#EFE5D7";
backColors[4] = "#5EB2D4";
backColors[8] = "#E1EFF8";
backColors[16] = "#3776C0";
backColors[64] = "#9ebbee";
backColors[32] = "#6bcae2";
backColors[128] = "white";

var colors = [];
colors[2] = "white";
colors[4] = "white";
colors[8] = "white";
colors[16] = "white";
colors[32] = "white";
colors[64] = "white";
colors[128] = "#2c3e50";

export default {
  name: "Chip",
  props: ["chip", "sizePx", "animationTimeMs"],
  data() {
    return {
      msg: "Welcome to Your Vue.js App"
    };
  },
  computed: {
    style() {
      return {
        fontSize: this.fontSizePx + "px",
        backgroundColor: this.backColor,
        color: this.color
      };
    },
    fontSizePx() {
      var n = Math.floor(Math.log(this.chip.value) / Math.log(10));
      var b = Math.floor(this.sizePx / 1.5);
      return b * (n < 8 ? fontSizeCoefs[n] : fontSizeCoefs[7]);
    },
    backColor: function() {
      return backColors[this.chip.value] || backColors[128];
    },
    color: function() {
      return colors[this.chip.value] || colors[128];
    }
  },
  watch: {
    "chip.value": function() {
      var el = this.$el;
      if (el) {
        var d = this.animationTimeMs + "ms";
        el.style["-webkit-animation"] = el.style.animation =
          "chip-value-changed " + d;
        el.style.transition = "background-color " + d;
        el.style["-webkit-transition"] = "-web-kit-background-color " + d;
      }
    }
  },
  methods: {
    enter: function(el, done) {
      var self = this;
      if (this.chip.prevRelPos) {
        var p = this.chip.prevRelPos;
        el.style["-webkit-transform"] = el.style.transform =
          "translate(" + p.left + "px," + p.top + "px)";
        requestAnimationFrame(function() {
          requestAnimationFrame(function() {
            el.style.transition = "transform " + self.animationTimeMs + "ms";
            el.style["-webkit-transition"] =
              "-webkit-transform " + self.animationTimeMs + "ms";
            el.style["-webkit-transform"] = el.style.transform = "";
          });
        });
      } else {
        el.style["-webkit-animation"] = el.style.animation =
          "chip-appear " + this.animationTimeMs + "ms";
      }
    }
  }
};
</script>
