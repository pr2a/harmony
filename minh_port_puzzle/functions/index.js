const functions = require('firebase-functions');

// const axios = require('axios');
// const functions = require('firebase-functions');
// var admin = require('firebase-admin');
// var serviceAccount = require('./keys/benchmark_account_key.json');
// const { getRandomWallet } = require('./keygen');

// admin.initializeApp({
//     credential: admin.credential.cert(serviceAccount),
//     databaseURL: 'https://benchmark-209420.firebaseio.com'
// });
// const LEADER_ADDRESS = `http://35.166.90.140:30000`;
// // const LEADER_ADDRESS = `127.0.0.1:30000`;

// const firestore = admin.firestore();

// function validateEmail(email) {
//     var re = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
//     return re.test(String(email).toLowerCase());
// }
// function anonymousEmail(email) {
//     var emailPattern = /^([a-zA-Z0-9._-]+)@([a-zA-Z0-9.-]+)\.([a-zA-Z]{2,4})$/;
//     const res = email.match(emailPattern);
//     return (
//         '...' + res[1].slice(-3) + '@' + '...' + res[2].slice(-3) + '.' + res[3]
//     );
// }

exports.reg = functions.https.onRequest(async (req, res) => {
  res.set('Access-Control-Allow-Origin', '*');
  res.set('Access-Control-Allow-Methods', 'GET, POST');

  try {
    res.json({
      account: '0x0000000000000000000000000000000000000000',
      email: 'ek@harmony.one',
      timestamp: Date.now()
    });
  } catch (err) {
    res.json({});
  }
});

exports.play = functions.https.onRequest(async (req, res) => {
  res.set('Access-Control-Allow-Origin', '*');
  res.set('Access-Control-Allow-Methods', 'GET, POST');

  // const email = req.query.email;
  // console.log(email);
  try {
    res.json({
      txid: '0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef',
      timestamp: Date.now(),
      token_change: 133
    });
  } catch (err) {
    res.json({});
  }
});

exports.finish = functions.https.onRequest(async (req, res) => {
  res.set('Access-Control-Allow-Origin', '*');
  res.set('Access-Control-Allow-Methods', 'GET, POST');

  try {
    res.json({
      level: 2,
      rewards: 20000000000000000000,
      timestamp: Date.now()
    });
  } catch (err) {
    res.json({});
  }
});
