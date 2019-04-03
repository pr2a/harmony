const { processEnter, processWinner, processResult } = require('./rest/client');
const express = require('express');
const cors = require('cors');
const path = require('path');
const app = express();
app.use(cors());

const PORT = 60000;

app.use(express.static(path.join(__dirname, 'build')));

app.get('/', function(req, res) {
  res.sendFile(path.join(__dirname, 'build', 'index.html'));
});

app.get('/enter/:key&:amount', (req, res) => {
  const key = req.params.key;
  const amount = req.params.amount;
  console.log('received ENTER request with', key, ' ', amount);
  processEnter(key, amount, res);
  // console.log('result:', r);
  // res.json(r);
});

app.get('/winner', (req, res) => {
  processWinner(res);
  // console.log('result:', r);
  // res.json(r);
});

app.get('/result/:key', (req, res) => {
  const key = req.params.key;
  console.log('received RESULT request with', key);
  processResult(key, res);
  // console.log('result:', r);
  // res.json(r);
});

app.get('/test', (req, res) => {
  res.json({ result: 'hello' });
});

module.exports = app;
// app.listen(PORT, err => {
//   if (err) throw err;
//   console.log(`ready at http://localhost:${PORT}`);
// });
