// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('@grpc/grpc-js');
var app_pb = require('./app_pb.js');

function serialize_function_RequestBody(arg) {
  if (!(arg instanceof app_pb.RequestBody)) {
    throw new Error('Expected argument of type function.RequestBody');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_function_RequestBody(buffer_arg) {
  return app_pb.RequestBody.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_function_ResponseBody(arg) {
  if (!(arg instanceof app_pb.ResponseBody)) {
    throw new Error('Expected argument of type function.ResponseBody');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_function_ResponseBody(buffer_arg) {
  return app_pb.ResponseBody.deserializeBinary(new Uint8Array(buffer_arg));
}


var gRPCFunctionService = exports.gRPCFunctionService = {
  gRPCFunctionHandler: {
    path: '/function.gRPCFunction/gRPCFunctionHandler',
    requestStream: false,
    responseStream: false,
    requestType: app_pb.RequestBody,
    responseType: app_pb.ResponseBody,
    requestSerialize: serialize_function_RequestBody,
    requestDeserialize: deserialize_function_RequestBody,
    responseSerialize: serialize_function_ResponseBody,
    responseDeserialize: deserialize_function_ResponseBody,
  },
};

exports.gRPCFunctionClient = grpc.makeGenericClientConstructor(gRPCFunctionService);
