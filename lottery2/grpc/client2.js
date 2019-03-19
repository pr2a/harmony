const PROTO_PATH = __dirname + '/message.proto';

const grpc = require('grpc');
const protoLoader = require('@grpc/proto-loader');
const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true
});
const messageProto = grpc.loadPackageDefinition(packageDefinition).message;

const process = msg =>
  function main() {
    let client = new messageProto.ClientService(
      'localhost:30000',
      grpc.credentials.createInsecure()
    );
    let user;
    if (process.argv.length >= 3) {
      user = process.argv[2];
    } else {
      user = 'world';
    }
    client.sayHello({ name: user }, function(err, response) {
      console.log('Greeting:', response.message);
    });
  };

main();

// // Process processes message.
// func(client * Client) Process(message * Message) * Response {
//     ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
//     defer cancel()
//     response, err := client.clientServiceClient.Process(ctx, message)
//     if err != nil {
//         log.Fatalf("Getting error when processing message: %s", err)
//     }
//     return response
// }
