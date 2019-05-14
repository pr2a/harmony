var admin = require('firebase-admin');

var serviceAccount = require('./keys/benchmark-firebase-db-key.json');

admin.initializeApp({
  credential: admin.credential.cert(serviceAccount),
  databaseURL: 'https://benchmark-209420.firebaseio.com'
});

const firestore = admin.firestore();

module.exports = {
  firestore
};
