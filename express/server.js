const {
  processEnter,
  processWinner,
  processResult,
  newSession
} = require('./client');
const express = require('express');
const next = require('next');
const cors = require('cors');
const firestore = require('./db');

const port = parseInt(process.env.PORT, 10) || 80;
const dev = process.env.NODE_ENV !== 'production';
const app = next({ dev });
const handle = app.getRequestHandler();

app.prepare().then(() => {
  const server = express();
  server.use(cors());

  server.get('/enter/:key&:amount&:email', (req, res) => {
    const key = req.params.key;
    const amount = req.params.amount;
    const email = req.params.email;
    console.log('received ENTER request with', key, ' ', amount, 'email:', email);
    processEnter(key, amount, email, res);
  });

  server.get('/winner/:admin', (req, res) => {
    const admin = req.params.admin;
    processWinner(admin, res);
  });

  server.get('/result', (req, res) => {
    console.log('received RESULT request');
    processResult(res);
  });

  server.get('/new', (req, res) => {
    console.log('received NEW request with');
    newSession(res);
  });
  server.get('/test', (req, res) => {
    res.json({ result: 'hello' });
  });

  server.get('*', (req, res) => {
    return handle(req, res);
  });

  server.listen(port, err => {
    if (err) throw err;
    console.log(`> Ready on http://localhost:${port}`);
  });
});
