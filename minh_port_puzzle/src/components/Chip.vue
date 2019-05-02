<template>
  <div class="chip" :style="style">{{ value }}</div>
</template>

<script>
var backColors = ["#00E0FF", "#482AFF", "#00AEE9", "#69FABD"];
var colors = ["#706B5E", "white", "white", "#19586D"];

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
