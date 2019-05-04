<template>
  <div class="chip" :style="style">{{ value }}</div>
</template>

<script>
const colors = [
  {
    bg: "#00AEE9",
    fg: "white"
  },
  {
    bg: "#00E0FF",
    fg: "#1B2A5E"
  },
  {
    bg: "#69FABD",
    fg: "#1B2A5E"
  },
  {
    bg: "#482AFF",
    fg: "white"
  },
  {
    bg: "#BDFF00",
    fg: "#1B295E"
  },
  {
    bg: "#758796",
    fg: "white"
  },
  {
    bg: "#C6F6EF",
    fg: "#1B295E"
  },
  {
    bg: "#19586D",
    fg: "white"
  },
  {
    bg: "#01DBA7",
    fg: "#1B295E"
  }
];

export default {
  name: "Chip",
  props: ["value"],
  computed: {
    style() {
      return {
        fontSize: "2.5em",
        backgroundColor: this.backColor,
        color: this.color
      };
    },
    backColor: function() {
      return colors[this.value % colors.length].bg;
    },
    color: function() {
      return colors[this.value % colors.length].fg;
    }
  }
};
</script>
