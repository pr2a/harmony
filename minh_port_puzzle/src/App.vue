<style lang="less">
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}
body,
html {
  font-family: "Source Sans Pro", Arial, sans-serif;
  color: #2c3e50;
  font-size: 18px;
  margin: 0;
  background-color: #fcf7eb;
}

.main-container {
  display: flex;
  flex-direction: column;
  margin: 0 auto;
  padding: 1em;
  /* Hack to improve transition performance on mobile devices. It enables GPU rendering. */
  transform: translateZ(0);
  -webkit-transform: translateZ(0);
}

.main-container .score-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.main-container .score-container .game-aim {
  font-weight: bold;
  font-size: 4em;
  text-align: center;
  color: white;
  background-color: #35495e;
  border-radius: 5% / 9%;
  width: 38%;
}

.main-container .score-container .game-aim-reached {
  text-shadow: 0 0 20px;
  animation: pulse 1s 3;
  -webkit-animation: pulse 1s 3;
  transition: text-shadow 3s;
  -webkit-transition: -webkit-text-shadow 3s;
}

@keyframes pulse {
  0% {
    transform: scale(1);
  }

  50% {
    transform: scale(1.1);
  }

  100% {
    transform: scale(1);
  }
}

@-webkit-keyframes pulse {
  0% {
    -webkit-transform: scale(1);
  }

  50% {
    -webkit-transform: scale(1.1);
  }

  100% {
    -webkit-transform: scale(1);
  }
}

.main-container .score-container .scores {
  font-weight: bold;
  display: flex;
  justify-content: space-between;
  width: 50%;
}

.main-container .score-container .scores .score {
  position: relative;
  text-align: center;
  color: white;
  background-color: #9aa4af;
  border-radius: 0.25em;
  width: 50%;
  padding-top: 2%;
}

.main-container .score-container .scores .score .score-inc {
  position: absolute;
  left: 0;
  color: #2c3e50;
  width: 100%;
  animation: up-disappear 1.5s;
  -webkit-animation: up-disappear 1.5s;
}

.main-container .score-container .scores .score .label {
  color: white;
  font-size: 1rem;
}

.main-container .game-controls {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  width: 100%;
}

.main-container .game-controls .button {
  background-color: #35495e;
  border: none;
  border-radius: 5% / 9%;
  width: 25%;
  height: 75%;
  color: white;
  outline: none;
  font-weight: bold;
  overflow: hidden;
  cursor: pointer;
  animation: appearing 0.5s;
  -webkit-animation: appearing 0.5s;
}

.main-container .game-controls .size-control {
  width: 75%;
  font-weight: bold;
  display: flex;
  justify-content: space-between;
  align-items: center;
  animation: appearing 0.5s;
  -webkit-animation: appearing 0.5s;
}

.main-container .game-controls .size-control input[type="radio"]:not(checked) {
  position: absolute;
  opacity: 0;
}

.main-container .game-controls .size-control input[type="radio"] + label {
  cursor: pointer;
  border: 5px solid #35495e;
  height: 30px;
  width: 30px;
  display: inline-flex;
  border-radius: 50%;
  justify-content: center;
  align-items: center;
  position: relative;
  vertical-align: middle;
  font-weight: bold;
  transition: background-color 0.3s, color 0.3s;
  -webkit-transition: -webkit-background-color 0.3s, -webkit-color 0.3s;
}

.main-container
  .game-controls
  .size-control
  input[type="radio"]:checked
  + label {
  background: #41b883;
  color: white;
}

.main-container .game-awards-container {
  display: flex;
  justify-content: space-around;
  align-items: center;
  margin-top: 2%;
}

.main-container .game-awards-container .award {
  border-radius: 15% / 50%;
  text-align: center;
  font-size: 1.2em;
}

.main-container .game-awards-container .award .like {
  width: 25%;
  top: 2px;
  position: relative;
}

.main-container .game-awards-container .award-not-obtained {
  color: #dee8ff;
  background-color: #9aa4af;
  fill: #dee8ff;
}

.main-container .game-awards-container .award-obtained {
  color: white;
  text-shadow: 0px 0px 4px orange;
  background: linear-gradient(to top right, gold, yellow);
  box-shadow: 0 2px 4px 0 black;
  fill: white;
  filter: drop-shadow(-1px -1px 0 orange);
}

.main-container .collect-all-awards {
  text-align: center;
  width: 100%;
  opacity: 0;
  margin-top: 4px;
}

.main-container .collect-all-awards span {
  border: 1px solid #2c3e50;
  border-radius: 7% / 50%;
  padding: 0 3px 0 3px;
}

.appearing07 {
  animation: appearing07 1s;
  -webkit-animation: appearing07 1s;
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
  background-color: #f96829;
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
  font-size: 1em;
  width: 4em;
  height: 2em;
  padding: 0;
  outline: none;
}

.logo {
  margin: 1em auto;
  color: red;
  height: 4em;
  width: 5em;
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

.flex-grow {
  flex: auto;
  min-width: 0;
  min-height: 0;
}
.text-right {
  text-align: right;
}

.info-item {
  background-color: #fff;
  border-radius: 0.5em;
  border: 0.15em solid #979797;
  padding: 1em;
  border-radius: 0.25em;
  padding: 0.5em;
  width: 6em;
  margin-left: 1em;
  text-align: center;
  .label {
    font-size: 0.5em;
    margin-bottom: 0.5em;
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
