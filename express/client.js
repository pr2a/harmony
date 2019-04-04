const axios = require('axios');
const { firestore } = require('./db');

const OWNER_KEY =
  '27978f895b11d9c737e1ab1623fde722c04b4f9ccb4ab776bf15932cc72d7c66';
// const LEADER_ADDRESS = '54.85.140.165';
const LEADER_ADDRESS = '18.236.187.250';
// const LEADER_ADDRESS = '127.0.0.1';

const PLAYERS = firestore.collection('players');
const CURRENT_SESSION = firestore.collection('lottery').doc('current_session');

const checkActiveSession = async () => {
  const currentSession = await CURRENT_SESSION.get();
  const data = currentSession.data();
  console.log('calling checkActiveSesion: ', data.active);
  return data.active;
};

const HasEntered = async player => {
  const currentSession = await CURRENT_SESSION.get();
  const data = currentSession.data();
  const players = await firestore
    .collection(`players`)
    .where('key', '==', player)
    .where('session', '==', data.session)
    .get();
  const empty = players.empty;
  return !empty;
};

const newSession = async res => {
  try {
    console.log('new session 1');
    const currentSession = await CURRENT_SESSION.get();
    const data = currentSession.data();
    if (data.active) {
      res.json({
        success: false,
        message: 'The current session is still active'
      });
      return;
    }
    await firestore
      .collection('lottery')
      .doc('current_session')
      .set({ session: data.session + 1, active: true });
    console.log('new session 2');
    return res.json({
      success: true,
      session: data.session + 1,
      message: `Session ${data.session + 1} is created`
    });
  } catch (err) {
    res.json({ success: false, message: 'Failed to call api' });
    return;
  }
};

const makeSessionOver = async () => {
  try {
    console.log('making session over');
    const currentSession = await CURRENT_SESSION.get();
    const data = currentSession.data();
    console.log(data);
    if (!data.active) {
      return true;
    }
    console.log('setting it not active');
    await firestore
      .collection('lottery')
      .doc('current_session')
      .set({ session: data.session, active: false });
    return true;
  } catch (err) {
    return false;
  }
};

const addNewPlayer = async (key, amount = 1) => {
  const currentSession = await CURRENT_SESSION.get();
  const data = currentSession.data();
  try {
    PLAYERS.add({
      key,
      amount,
      session: data.session
    });
    return true;
  } catch (err) {
    return false;
  }
};

const processEnter = async (key, amount, res) => {
  try {
    const active = await checkActiveSession();
    if (!active) {
      res.json({
        success: false,
        message: 'The current session is over'
      });
      return;
    }
    // if (await HasEntered(key)) {
    //   res.json({
    //     success: false,
    //     message: `Player ${key} had entered before`
    //   });
    //   return;
    // }
    const { data } = await axios.get(
      `http://${LEADER_ADDRESS}:30000/enter?key=${key}&amount=${amount}`
    );
    console.log(data);
    if (data.success) {
      if (await addNewPlayer(key, amount)) {
        res.json({ success: true, message: `Player ${key} entered` });
      } else {
        res.json({
          success: false,
          message: `Failed to add new player into db`
        });
      }
    } else {
      res.json({
        success: false,
        message: 'Failed to process enter in blockchain.'
      });
    }
  } catch (err) {
    res.json({ success: false, message: 'Failed to call api' });
    console.log(err);
    return;
  }
};

const processWinner = async (adminKey, res) => {
  try {
    console.log('checking active session or not');
    const active = await checkActiveSession();
    if (!active) {
      console.log('go here', active);
      res.json({
        success: false,
        message: 'The current session is over'
      });
      return;
    }
    console.log('go here2', active);
    console.log('the  session is active');
    if (adminKey !== OWNER_KEY) {
      res.json({
        success: false,
        message: 'Only owner key has the right to process PickWinner'
      });
      return;
    }
    const { data } = await axios.get(`http://${LEADER_ADDRESS}:30000/winner`);
    console.log('result of calling /winner', data);
    if (data.success) {
      if (await makeSessionOver()) {
        res.json({ success: true, message: 'A winner is picked.' });
      } else {
        res.json({ success: false, message: 'Failed to set the session over' });
      }
    } else {
      res.json({
        success: false,
        message: 'Failed to process pickWinner in blockchain.'
      });
    }
  } catch (err) {
    res.json({ success: false, message: 'Failed to call api' });
    console.log(err);
    return;
  }
};

const processResult = async res => {
  try {
    const { data } = await axios.get(
      `http://${LEADER_ADDRESS}:30000/result?key=${OWNER_KEY}`
    );

    const currentSession = await CURRENT_SESSION.get();
    const sessionData = currentSession.data();

    console.log('sessionData', sessionData);
    let ret = {
      players: data.players,
      balances: data.balances,
      session: sessionData.session,
      success: true
    };
    ret.message = `There are ${
      ret.players.length
    } players participating in the lottery session ${ret.session}`;
    res.json(ret);
  } catch (err) {
    res.json({ success: false, message: 'Failed to call api' });
    console.log(err);
    return;
  }
};

module.exports = {
  processEnter,
  processWinner,
  processResult,
  newSession
};
