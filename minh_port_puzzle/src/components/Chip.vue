<template>
  <!-- <transition :css="false" @enter="enter"> -->
  <div class="chip" :style="style">{{ value }}</div>
  <!-- </transition> -->
</template>

<script>
var fontSizeCoefs = [1, 1, 0.8, 0.65, 0.5, 0.4, 0.35, 0.32];
var backColors = ["#EFE5D7", "#5EB2D4", "#E1EFF8", "#3776C0", "#9ebbee"];

var colors = ["#706B5E", "white", "#19586D", "white", "#2f3747"];

export default {
  name: "Chip",
  props: ["value", "sizePx", "animationTimeMs"],
  data() {
    return {};
  },
  computed: {
    style() {
      return {
        fontSize: "30px",
        backgroundColor: this.backColor,
        color: this.color
      };
    },
    backColor: function() {
      return backColors[this.value % backColors.length];
    },
    color: function() {
      return colors[this.value % colors.length];
    }
  },
  watch: {
    value: function() {
      var el = this.$el;
      if (el) {
        var d = this.animationTimeMs + "ms";
        el.style["-webkit-animation"] = el.style.animation =
          "chip-changed " + d;
        el.style.transition = "background-color " + d;
        el.style["-webkit-transition"] = "-web-kit-background-color " + d;
      }
    }
  },
  methods: {
    enter: function(el, done) {
      var self = this;
    }
  }
};
</script>
