const axios = require("axios");

// 1+3 = 4 leaders
const API_URL = ["http://127.0.0.1:30000"];

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
async function register(email) {
  var res;
  if (validateEmail(email)) {
    res = await axios.post(`${API_URL}/v1/reg`, {
      id: email
    });
  }
  return res;
}

async function play(key, stake) {
  var res;
  if (validateStake(stake) && validatePrivateKey(key)) {
    res = await axios.post(`${API_URL}/v1/play`, {
      key: key,
      stake: stake
    });
  }
  return res;
}

async function finish(key, txid, level, seq) {
  var res;
  if (validatePrivateKey(key)) {
    res = await axios.post(`${API_URL}/v1/finish`, {
      key: key,
      txid: txid,
      level: level,
      seq: seq
    });
  }
  return res;
}

register("chao@harmony.one")
  .then(res => {
    if (res && res.data) {
      console.log(res.data);
    }
  })
  .catch(err => console.log(err.response));
play("0xabc", 23)
  .then(res => {
    if (res && res.data) {
      console.log(res.data);
    }
  })
  .catch(err => console.log(err.response));
finish("0xabc", "0x2341", 22, "ULDRU")
  .then(res => {
    if (res && res.data) {
      console.log(res.data);
    }
  })
  .catch(err => console.log(err.response));
