const axios = require('axios');

const processEnter = async (key, amount, res) => {
  try {
    const { data } = await axios.get(
      `http://localhost:30000/enter?key=${key}&amount=${amount}`
    );
    console.log(data);
    if (data.success) {
      res.json({ success: true });
    } else {
      res.json({ success: false });
    }
  } catch (err) {
    res.json({ success: false, message: 'failed to call api' });
    console.log(err);
    return;
  }
};

const processWinner = async res => {
  try {
    const data = await axios.get(`http://localhost:30000/winner`);
    if (data.success) {
      res.json({ success: true });
    } else {
      res.json({ success: false });
    }
  } catch (err) {
    res.json({ success: false, message: 'failed to call api' });
    console.log(err);
    return;
  }
};

const processResult = async (key, res) => {
  try {
    const data = await axios.get(`http://localhost:30000/result?key=${key}`);
    res.json(data.data);
  } catch (err) {
    res.json({ success: false, message: 'failed to call api' });
    console.log(err);
    return;
  }
};

module.exports = {
  processEnter,
  processWinner,
  processResult
};
