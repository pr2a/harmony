const functions = require('firebase-functions');
// The Firebase Admin SDK to access the Firebase Realtime Database.
const admin = require('firebase-admin');
admin.initializeApp();

const app = require('./lottery2/server');

exports.helloWorld = functions.https.onRequest(app);
