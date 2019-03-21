const axios = require('axios');
const msg = require('./message_pb');
const services = require('./message_grpc_pb');
const grpc = require('grpc');

const API = axios.create({ baseURL: 'localhost:30000/' });

// const PROXY_CONFIG = {
//   headers: { 'Access-Control-Allow-Origin': '*' }
// };

// const response = await dict.get(`/api/webster/similar/${term}`);

// curl 'http://127.0.0.1:31313/enter?key=27978f895b11d9c737e1ab1623fde722c04b4f9ccb4ab776bf15932cc72d7c66&amount=1'
// curl 'http://127.0.0.1:31313/enter?key=371cb68abe6a6101ac88603fc847e0c013a834253acee5315884d2c4e387ebca&amount=1'
// curl 'http://127.0.0.1:31313/enter?key=3f8af52063c6648be37d4b33559f784feb16d8e5ffaccf082b3657ea35b05977&amount=1'
// curl 'http://127.0.0.1:31313/enter?key=df77927961152e6a080ac299e7af2135fc0fb02eb044d0d7bbb1e8c5ad523809&amount=1'

// curl 'http://127.0.0.1:31313/result?key=27978f895b11d9c737e1ab1623fde722c04b4f9ccb4ab776bf15932cc72d7c66'
// curl 'http://127.0.0.1:31313/winner'

const client = new services.ClientServiceClient(
  'localhost:30000',
  grpc.credentials.createInsecure()
);

const processEnter = async (key, amount, res) => {
  // let message = createMessage('ENTER', key, amount);
  // const result = await client.process(message);
  // const result = await client.process(message);
  try {
    const result = await API.get(`/enter?key=${key}&amount=${amount}`);
  } catch (err) {
    res.json({ success: false, message: 'failed to call api' });
    return;
  }

  console.log(result);
  if (result.success && result.success === true) {
    res.json({ success: true });
  } else {
    res.json({ success: false });
  }
};

const processWinner = async res => {
  // let message = createMessage('WINNER');

  try {
    const result = await API.get(`/winner`);
  } catch (err) {
    res.json({ success: false, message: 'failed to call api' });
    return;
  }
  if (result.success && result.success === true) {
    res.json({ success: true });
  } else {
    res.json({ success: false });
  }
};

const processResult = async (key, res) => {
  // let message = createMessage('RESULT', key);
  // client.process(message, function(err, response) {
  //   console.log('Response of RESULT:', response);
  //   res.json(response);
  // });
  try {
    const result = await API.get(`/result?key=${key}`);
  } catch (err) {
    res.json({ success: false, message: 'failed to call api' });
    return;
  }
  console.log(result);
  res.json(result);
};

const createMessage = (requestType, key = null, amount = null) => {
  let message = new msg.Message();
  let lotteryRequest = new msg.LotteryRequest();
  switch (requestType) {
    case 'ENTER':
      lotteryRequest.setType(msg.LotteryRequest.Type.ENTER);
      break;
    case 'RESULT':
      lotteryRequest.setType(msg.LotteryRequest.Type.RESULT);
      break;
    case 'WINNER':
      lotteryRequest.setType(msg.LotteryRequest.Type.PICK_WINNER);
      break;
  }
  if (key) {
    lotteryRequest.setPrivateKey(key);
  }
  if (amount) {
    lotteryRequest.setAmount(amount);
  }

  message.setReceiverType(msg.ReceiverType.CLIENT);
  message.setServiceType(msg.ServiceType.CLIENT_SUPPORT);
  message.setType(msg.MessageType.LOTTERY_REQUEST);
  message.setLotteryRequest(lotteryRequest);
  return message;
};

module.exports = {
  processEnter,
  processWinner,
  processResult
};
