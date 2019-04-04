const admin = require('firebase-admin');
// const serviceAccount = require('./keys/serviceAccountKey.json');
const serviceAccount = require('./keys/benchmark_account_key.json');

admin.initializeApp({
  credential: admin.credential.cert(serviceAccount),
  databaseURL: 'https://benchmark-209420.firebaseio.com'
});

const firestore = admin.firestore();

module.exports = {
  firestore
};

// admin.initializeApp({
//   credential: admin.credential.cert(serviceAccount),
//   databaseURL: 'https://benchmark-209420.firebaseio.com'
// });
