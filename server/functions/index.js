const { getRandomWallet } = require('./keygen');
const functions = require('firebase-functions');
// const { firestore } = require('./db');

var admin = require('firebase-admin');

var serviceAccount = require('./keys/benchmark_account_key.json');

admin.initializeApp({
  credential: admin.credential.cert(serviceAccount),
  databaseURL: 'https://benchmark-209420.firebaseio.com'
});

const firestore = admin.firestore();

// Create and Deploy Your First Cloud Functions
// https://firebase.google.com/docs/functions/write-firebase-functions

const OWNER_KEY =
  '27978f895b11d9c737e1ab1623fde722c04b4f9ccb4ab776bf15932cc72d7c66';
// const LEADER_ADDRESS = '54.85.140.165';
// const LEADER_ADDRESS = '54.175.2.80';
// const PLAYERS = firestore.collection('players');
// const CURRENT_SESSION = firestore.collection('session').doc('current_session');

// const checkActiveSession = async () => {
//   const currentSession = await CURRENT_SESSION.get();
//   const data = currentSession.data();
//   console.log('calling checkActiveSesion: ', data.active);
//   return data.active;
// };

// const HasEntered = async player => {
//   const currentSession = await CURRENT_SESSION.get();
//   const data = currentSession.data();
//   const players = await firestore
//     .collection(`players`)
//     .where('key', '==', player)
//     .where('session', '==', data.session)
//     .get();
//   const empty = players.empty;
//   return !empty;
// };

// console.log('start');
// var citiesRef = db.collection('cities');
// try {
//   var allCitiesSnapShot = await citiesRef.get();
//   allCitiesSnapShot.forEach(doc => {
//     console.log(doc.id, '=>', doc.data().name);
//   });
//   console.log('end');
// } catch (err) {
//   console.log('Error getting documents', err);
// }

exports.enter = functions.https.onRequest(async (req, res) => {
  const email = req.query.email;
  try {
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
        const { address, private_key } = getRandomWallet();

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

// const _ = require('lodash');
// const FuzzyTrie = require('fuzzytrie');

// const MAX_DISTANCE = 1;
// const MAX_RETURN = 12;

// let trie = new FuzzyTrie();
// const jsonData = require('./en_vn.json');
// let dict = {};
// _.forEach(jsonData, item => {
//     const word = item.word;
//     if (!word.includes(' ')) {
//         trie.add(word);
//         dict[item.word] = item;
//     }
// });

// const functions = require('firebase-functions');

// // English to Vietnamese.
// exports.envn = functions.https.onRequest((req, res) => {
//     const word = req.query.word;
//     if (word in dict) {
//         res.json({ def: dict[word] });
//     } else {
//         res.json({});
//     }
// });

// exports.similar = functions.https.onRequest((req, res) => {
//     const word = req.query.word;
//     let result = [];
//     _.forEach(trie.find(word, MAX_DISTANCE), (value, key) =>
//         result.push({ word: key, value })
//     );

//     result = _.sortBy(result, item => item.value);
//     res.json({ similar: result.slice(0, MAX_RETURN) });
// });
