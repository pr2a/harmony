const functions = require('firebase-functions');

exports.reg = functions.https.onRequest(async (req, res) => {
  res.set('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Methods', 'POST,GET,OPTIONS,PUT,DELETE');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type,Accept');

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
  res.setHeader('Access-Control-Allow-Methods', 'POST,GET,OPTIONS,PUT,DELETE');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type,Accept');

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
  res.setHeader('Access-Control-Allow-Methods', 'POST,GET,OPTIONS,PUT,DELETE');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type,Accept');

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
