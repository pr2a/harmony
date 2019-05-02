const axios = require("axios");
const pinfo = require("./playinfo");
const API_URL = "http://127.0.0.1:30000/v1";

function validateEmail(email) {
  var re = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
  return re.test(String(email).toLowerCase());
}
function validateStake(stake) {
  if (typeof stake === "number" && stake > 0) return true;
  if (parseFloat(stake) && parseFloat(stake) > 0) return true;
  return false;
}

function validatePrivateKey(key) {
  return true;
}

function sendReq(url, data) {
  return axios.post(API_URL + url, data, {});
}

function register(email) {
  if (validateEmail(email)) {
    sendReq("/reg", {
      id: email
    })
      .then(res => {
        if (res && res.data) {
          console.log(res.data);
        }
      })
      .catch(err => {
        console.log(err.response);
      });
  }
}

function play(key, stake) {
  if (validateStake(stake) && validatePrivateKey(key)) {
    sendReq("/play", {
      key: key,
      stake: stake
    })
      .then(res => {
        if (res && res.data) {
          console.log(res.data);
        }
      })
      .catch(err => {
        console.log(err.response);
      });
  }
}

function finish(key, txid, playinfo) {
  if (validatePrivateKey(key)) {
    sendReq("/finish", {
      key: key,
      txid: txid,
      playinfo: playinfo
    })
      .then(res => {
        if (res && res.data) {
          console.log(res.data);
        }
      })
      .catch(err => {
        console.log(err.response);
      });
  }
}

register("chao@harmony.one");
play("0xabc", 23);
finish("0xabc", "0x2341", 22, pinfo.playHistory);
