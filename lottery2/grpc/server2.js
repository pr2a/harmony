var PROTO_PATH = __dirname + '/message.proto';

var grpc = require('grpc');
var protoLoader = require('@grpc/proto-loader');
var packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true
});
var messageProto = grpc.loadPackageDefinition(packageDefinition).message;

/**
 * Implements the SayHello RPC method.
 */
function process(call, callback) {
  callback(null, {
    receiver_type: 'NEWNODE',
    service_type: 'CONSENSUS',
    message_type: 'COMMITTED',
    response: {
      lottery_response: {
        players: ['minh', 'mark'],
        balances: [1, 2]
      }
    }
  });
}

/**
 * Starts an RPC server that receives requests for the Greeter service at the
 * sample server port
 */
function main() {
  var server = new grpc.Server();
  server.addService(messageProto.ClientService.service, { process });
  server.bind('0.0.0.0:50051', grpc.ServerCredentials.createInsecure());
  console.log('start');
  server.start();
}

main();

protoc --proto_path = src--js_out = library = whizz / ponycopter, binary: build / gen src / foo.proto src / bar / baz.proto
protoc --proto_path=./ --js_out=library=, binary: build / gen src / foo.proto src / bar / baz.proto
