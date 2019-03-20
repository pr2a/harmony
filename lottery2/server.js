const {
  processEnter,
  processWinner,
  processResult
} = require('./grpc/client2');
const express = require('express');
const app = express();

const PORT = 60000;

app.get('/enter/:key&:amount', (req, res) => {
  const key = req.params.key;
  const amount = req.params.amount;
  console.log('received ENTER request with', key, ' ', amount);
  const r = processEnter(key, amount);

  console.log('result:', r);
  res.json(r);
});

app.get('/winner', (req, res) => {
  console.log('received WINNER request with');
  const r = processWinner();
  console.log('result:', r);
  res.json(r);
});

app.get('/result/:key', (req, res) => {
  const key = req.params.key;
  console.log('received RESULT request with', key);
  const r = processResult(key);
  console.log('result:', r);
  res.json(r);
});

app.listen(PORT, err => {
  if (err) throw err;
  console.log(`ready at http://localhost:${PORT}`);
});
