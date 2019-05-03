<style lang="less">
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}
html {
  background-image: linear-gradient(
    180deg,
    #00aee9 0%,
    rgba(105, 250, 189, 0) 119.53%
  );
  height: 100%;
}
body,
#app {
  height: 100%;
}
body,
html {
  font-family: "Source Sans Pro", Arial, sans-serif;
  color: #2c3e50;
  font-size: 18px;
  margin: 0;
}

.main-container {
  display: flex;
  flex-direction: column;
  margin: 0 auto;
  padding: 1em;
  /* Hack to improve transition performance on mobile devices. It enables GPU rendering. */
  transform: translateZ(0);
}

.main-container .score-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

@keyframes up-disappear {
  0% {
    opacity: 0.7;
  }

  100% {
    opacity: 0;
    transform: translateY(-40px);
  }
}

@-webkit-keyframes up-disappear {
  0% {
    opacity: 0.7;
  }

  100% {
    opacity: 0;
    -webkit-transform: translateY(-40px);
  }
}

@keyframes chip-value-changed {
  0% {
    transform: scale(1);
  }

  50% {
    transform: scale(1.2);
  }

  100% {
    transform: scale(1);
  }
}

.btn-primary,
.btn-mini {
  background-color: #482bff;
  border-radius: 0.3em;
  border: 0;
  color: #fff;
  outline: none;
  &:disabled {
    opacity: 0.5;
    background-color: #aaa;
  }
}

.btn-primary {
  padding: 0.5em 1.5em;
  font-size: 1.2em;
}

.btn-mini {
  display: inline-block;
  font-size: 0.8em;
  width: 2em;
  height: 2em;
  padding: 0;
  outline: none;
}

.logo {
  margin: 1em auto;
  color: red;
  height: 2em;
  width: 2em;
  background-image: url(./assets/logo.svg);
  background-size: contain;
  background-repeat: no-repeat;
}

.flex-veritcal {
  display: flex;
  flex-direction: column;
}

.flex-horizontal {
  display: flex;
  flex-direction: row;
  align-items: center;
}
.flex-hv-center {
  display: flex;
  align-items: center;
  justify-content: center;
}

.flex-grow {
  flex: auto;
  min-width: 0;
  min-height: 0;
}
.text-right {
  text-align: right;
}

.info-item {
  .flex-horizontal;
  color: #fff;
  padding: 0.5em;
  width: 5em;
  margin-left: 1em;
  text-align: center;
  .content {
    margin-left: 0.5em;
    font-weight: bold;
  }
}

p {
  margin: 0;
}
p + p {
  margin-top: 1em;
}
</style>

<template>
  <div id="app">
    <router-view/>
  </div>
</template>

<script>
export default {
  name: "App"
};
</script>

<style>
</style>
