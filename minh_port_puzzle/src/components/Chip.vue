<template>
  <div class="chip" :style="style">{{ value }}</div>
</template>

<script>
var backColors = ["#EFE5D7", "#5EB2D4", "#E1EFF8", "#3776C0", "#9ebbee"];
var colors = ["#706B5E", "white", "#19586D", "white", "#2f3747"];

export default {
  name: "Chip",
  props: ["value"],
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
  }
};
</script>
