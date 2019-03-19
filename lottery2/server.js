const express = require('express');
const app = express();
const PORT = 9000;

app.get('/enter/:key&:amount', (req, res) => {
  const key = req.params.key;
  const amount = req.params.amount;
  console.log(key, ' ', amount);
  res.json({});
});

app.get('/winner', (req, res) => {
  res.json({});
});

app.get('/result/:key', (req, res) => {
  const key = req.params.key;
  console.log(key);
  res.json({});
});

app.listen(PORT, err => {
  if (err) throw err;
  console.log(`ready at http://localhost:${PORT}`);
});
