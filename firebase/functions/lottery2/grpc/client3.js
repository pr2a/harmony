const axios = require('axios');

const processEnter = async (key, amount, res) => {
  try {
    const result = await axios.get(
      `http://localhost:30000/enter?key=${key}&amount=${amount}`
    );
    console.log(result);
    if (result.success && result.success === true) {
      res.json({ success: true });
    } else {
      res.json({ success: false });
    }
  } catch (err) {
    res.json({ success: false, message: 'failed to call api' });
    return;
  }
};

const processWinner = async res => {
  try {
    const result = await axios.get(`http://localhost:30000/winner`);
    if (result.success && result.success === true) {
      res.json({ success: true });
    } else {
      res.json({ success: false });
    }
  } catch (err) {
    res.json({ success: false, message: 'failed to call api' });
    return;
  }
};

const processResult = async (key, res) => {
  try {
    const result = await axios.get(`http://localhost:30000/result?key=${key}`);
    console.log(result);
    res.json(result);
  } catch (err) {
    res.json({ success: false, message: 'failed to call api' });
    return;
  }
};

module.exports = {
  processEnter,
  processWinner,
  processResult
};
