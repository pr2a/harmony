const msg = require('./message_pb');
const services = require('./message_grpc_pb');
const grpc = require('grpc');

const client = new services.ClientServiceClient(
  'localhost:30000',
  grpc.credentials.createInsecure()
);

const processEnter = (key, amount) => {
  let message = createMessage('ENTER', key, amount);

  client.process(message, function(err, response) {
    console.log('Greeting:', response);
  });
};

const processWinner = () => {
  let message = createMessage('WINNER');

  client.process(message, function(err, response) {
    console.log('Greeting:', response);
  });
};

const processResult = key => {
  let message = createMessage('RESULT', key);

  client.process(message, function(err, response) {
    console.log('Greeting:', response);
  });
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
