const { firestore } = require('./db');

firestore
  .collection('lottery')
  .doc('current_session')
  .set({
    session: 0
  })
  .then(ref => console.log('added', ref.id))
  .catch(err => console.log('error:', err));
