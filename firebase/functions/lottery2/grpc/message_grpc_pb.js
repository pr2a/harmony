// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var message_pb = require('./message_pb.js');

function serialize_message_Message(arg) {
  if (!(arg instanceof message_pb.Message)) {
    throw new Error('Expected argument of type message.Message');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_message_Message(buffer_arg) {
  return message_pb.Message.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_message_Response(arg) {
  if (!(arg instanceof message_pb.Response)) {
    throw new Error('Expected argument of type message.Response');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_message_Response(buffer_arg) {
  return message_pb.Response.deserializeBinary(new Uint8Array(buffer_arg));
}


// Client is the service used for any client-facing requests.
var ClientServiceService = exports.ClientServiceService = {
  process: {
    path: '/message.ClientService/Process',
    requestStream: false,
    responseStream: false,
    requestType: message_pb.Message,
    responseType: message_pb.Response,
    requestSerialize: serialize_message_Message,
    requestDeserialize: deserialize_message_Message,
    responseSerialize: serialize_message_Response,
    responseDeserialize: deserialize_message_Response,
  },
};

exports.ClientServiceClient = grpc.makeGenericClientConstructor(ClientServiceService);
