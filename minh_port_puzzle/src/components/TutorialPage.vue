<style scoped lang="less">
.content {
  margin: 3em;
}
.key {
  border-radius: 0.5em;
  border: 0.15em solid #979797;
  padding: 1em;
  background-color: #fff;
  display: block;
  width: 100%;
  overflow: auto;
}
.msg {
  padding: 1em;
  border-radius: 0.5em;
  color: #59504d;
  font-family: Fira Sans, sans-serif;
  border: 0.15em solid #979797;
  height: 10em;
  overflow: auto;
  background-color: #fff;
  li {
    list-style: none;
    &:before {
      content: "\1F449";
      margin-right: 0.5em;
    }
  }
  li + li {
    margin-top: 0.5em;
  }
}

.btn-primary {
  display: block;
  margin: auto;
}
</style>

<template >
  <div class="tutorial-page">
    <div class="content">
      <div class="logo"></div>
      <div class="msg">
        <ul>
          <li>Swipe or use arrow keys to move between chips</li>
          <li>The number on chip increments by one when you move onto it</li>
          <li>You complete the level once all the chips have the same number</li>
        </ul>
      </div>
    </div>
    <button class="btn-primary" @click="$emit('start')">Got it</button>
  </div>
</template>

<script>
export default {
  name: "TutorialPage"
};
</script>
