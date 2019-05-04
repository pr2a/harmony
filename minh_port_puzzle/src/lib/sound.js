const moveSound = require('../assets/move.wav');
const beginSound = require('../assets/begin.wav');
const endSound = require('../assets/end.wav');
const backgroundMusic = require('../assets/cryptic.ogg');
const backgroundMusicAudio = new Audio(backgroundMusic);

playSound = sound => {
  var audio = new Audio(sound);
  audio.play();
};

playAudio = audio => {
  audio.play();
};

playAudioLoop =  audio => {
  audio.play();
  audio.loop = true;
}

stopAudio = audio => {
  audio.pause();
  audio.currentTime = 0;
};

stopSound = audio => {
  audio.pause();
  audio.currentTime = 0;
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

playBackgroundMusic = () => {
  playAudioLoop(backgroundMusicAudio);
};

stopBackgroundMusic = () => {
  stopSound(backgroundMusicAudio);
};

playPostGameMusic = () => {
  playAudio(postGameMusicAudio);
};

module.exports = {
  playMoveSound,
  playBeginSound,
  playEndSound,
  playBackgroundMusic
};
