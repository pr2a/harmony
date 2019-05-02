const moveSound = require('../assets/move.wav');
const beginSound = require('../assets/begin.wav');
const endSound = require('../assets/end.wav');

playSound = sound => {
  var audio = new Audio(sound);
  audio.play();
};

playMoveSound = () => {
  playSound(moveSound);
};

playBeginSound = () => {
  playSound(beginSound);
};

playEndSound = () => {
  playSound(endSound);
};

module.exports = {
  playMoveSound,
  playBeginSound,
  playEndSound
};
