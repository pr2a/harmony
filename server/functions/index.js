const { getRandomWallet } = require('./keygen');
const functions = require('firebase-functions');
var admin = require('firebase-admin');
var serviceAccount = require('./keys/benchmark_account_key.json');

admin.initializeApp({
  credential: admin.credential.cert(serviceAccount),
  databaseURL: 'https://benchmark-209420.firebaseio.com'
});

// const LEADER_ADDRESS = `34.222.210.98`;
const LEADER_ADDRESS = `127.0.0.1`;

const firestore = admin.firestore();

function validateEmail(email) {
  var re = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
  return re.test(String(email).toLowerCase());
}

exports.existed = functions.https.onRequest(async (req, res) => {
  res.set('Access-Control-Allow-Origin', '*');
  res.set('Access-Control-Allow-Methods', 'GET, POST');

  const email = req.query.email;
  try {
    if (!validateEmail(email)) {
      res.json({ status: 'failed', message: 'invalid email' });
      return;
    }
    const active_session = await firestore
      .collection('session')
      .where('is_current', '==', true)
      .get();

    if (active_session.empty) {
      res.json({});
    } else {
      const existed = await firestore
        .collection('players')
        .where('email', '==', email)
        .get();

      if (existed.empty) {
        res.json({ status: 'success', joined: false });
      } else {
        let done = false;
        existed.forEach(async doc => {
          if (done) {
            return;
          }
          done = true;
          const data = doc.data();
          console.log(data);
          address = data.address;
          private_key = data.private_key;
          res.json({
            status: 'success',
            joined: true,
            address,
            private_key
          });
        });
      }
    }
  } catch (err) {
    res.json({});
  }
});

exports.enter = functions.https.onRequest(async (req, res) => {
  res.set('Access-Control-Allow-Origin', '*');
  res.set('Access-Control-Allow-Methods', 'GET, POST');

  const email = req.query.email;
  const private_key = req.query.private_key;
  const address = req.query.address;
  const funded = req.query.funded;
  try {
    if (!validateEmail(email)) {
      res.json({ status: 'failed', message: 'invalid email' });
      return;
    }
    const active_session = await firestore
      .collection('session')
      .where('is_current', '==', true)
      .get();

    active_session.forEach(async doc => {
      const session_id = doc.data().session_id;
      console.log('come 6', session_id);
      const existed = await firestore
        .collection('players')
        .where('session_id', '==', session_id)
        .where('email', '==', email)
        .get();

      if (existed.empty) {
        if (!funded) {
          const { data } = await axios.get(
            `LEADER_ADDRESS/fundme?address=${address}`
          );
          console.log('fundme:', data);
          if (!data || !data.success) {
            res.json({
              status: false,
              message: `Unable to fund your account`
            });
            return;
          }
        }

        const { data } = await axios.get(
          `http://${LEADER_ADDRESS}:30000/enter?key=${key}&amount=1`
        );
        console.log(data);
        if (data && data.success) {
          await firestore.collection('players').add({
            email,
            private_key,
            address,
            keys_notified: false,
            result_notified: false,
            session_id
          });
          res.json({
            status: 'success',
            message: 'You have entered to the current lottery session.'
          });
        } else {
          res.json({
            status: false,
            message: 'Failed to process enter in blockchain.'
          });
          return;
        }
      } else {
        res.json({
          status: 'failed',
          message: 'Your email has been used in this session'
        });
      }
    });
  } catch (err) {
    res.json({});
  }
});

exports.current_session = functions.https.onRequest(async (req, res) => {
  try {
    const active_session = await firestore
      .collection('session')
      .where('is_current', '==', true)
      .get();

    if (active_session.empty) {
      res.json({});
    } else {
      active_session.forEach(async doc => {
        const data = doc.data();
        res.json({ deadline: data.deadline, session_id: data.session_id });
      });
    }
  } catch (err) {
    console.log(err);
    res.json({});
  }
});

exports.current_players = functions.https.onRequest(async (req, res) => {
  try {
    const active_session = await firestore
      .collection('session')
      .where('is_current', '==', true)
      .get();

    if (active_session.empty) {
      res.json({});
    } else {
      active_session.forEach(async doc => {
        const data = doc.data();
        const players = await firestore
          .collection('players')
          .where('session_id', '==', data.session_id)
          .get();
        let result = [];
        players.forEach(player => {
          result.push(player.data().address);
        });
        res.json({ current_players: result });
      });
    }
  } catch (err) {
    console.log(err);
    res.json({});
  }
});

exports.previous_winners = functions.https.onRequest(async (req, res) => {
  try {
    const winners = await firestore.collection('winners').get();
    let result = [];
    winners.forEach(winner => {
      result.push({ ...winner.data() });
    });
    res.json({ previous_winners: result });
  } catch (err) {
    console.log(err);
    res.json({});
  }
});
