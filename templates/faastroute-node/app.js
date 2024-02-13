const express = require('express')
const app = express()
const grpc = require('@grpc/grpc-js')
const redis = require('redis')
const moment = require('moment')
const index = require('./func')
const messages = require('./app_pb')
const services = require('./app_grpc_pb')

if ("LOGGING_NAME" in process.env) {
  const client = redis.createClient({url: process.env.LOGGING_URL, password: process.env.LOGGING_PASSWORD});
}

app.get('/', (req, res) => {
  if (req.header("Content-Type") == "application/json") {
    var workflow_id = Math.floor(Math.random() * 10000000).toString();
    if ("LOGGING_NAME" in process.env) {
      client.append(process.env.LOGGING_NAME, "EXECUTION_GET_JSON_START," + workflow_id + ",NA,1," + moment().format('MMMM Do YYYY, h:mm:ss a') + "\n");
    }
    var [reply, code] = func.function_handler({"request": req.body, "workflow_id": workflow_id, "request_type": "GET", "is_json": true});
    if ("LOGGING_NAME" in process.env) {
      client.append(process.env.LOGGING_NAME, "EXECUTION_GET_JSON_END," + workflow_id + ",NA,1," + moment().format('MMMM Do YYYY, h:mm:ss a') + "\n");
    }
    res.send(reply)
  }
  else {
    var workflow_id = Math.floor(Math.random() * 10000000).toString();
    if ("LOGGING_NAME" in process.env) {
      client.append(process.env.LOGGING_NAME, "EXECUTION_GET_TEXT_START," + workflow_id + ",NA,1," + moment().format('MMMM Do YYYY, h:mm:ss a') + "\n");
    }
    var [reply, code] = func.function_handler({"request": req.body, "workflow_id": workflow_id, "request_type": "GET", "is_json": false});
    if ("LOGGING_NAME" in process.env) {
      client.append(process.env.LOGGING_NAME, "EXECUTION_GET_TEXT_END," + workflow_id + ",NA,1," + moment().format('MMMM Do YYYY, h:mm:ss a') + "\n");
    }
    res.send(reply)
  }
})

app.post('/', (req, res) => {
  if (req.header("Content-Type") == "application/json") {
    var workflow_id = Math.floor(Math.random() * 10000000).toString();
    if ("LOGGING_NAME" in process.env) {
      client.append(process.env.LOGGING_NAME, "EXECUTION_POST_JSON_START," + workflow_id + ",NA,1," + moment().format('MMMM Do YYYY, h:mm:ss a') + "\n");
    }
    var [reply, code] = func.function_handler({"request": req.body, "workflow_id": workflow_id, "request_type": "POST", "is_json": true});
    if ("LOGGING_NAME" in process.env) {
      client.append(process.env.LOGGING_NAME, "EXECUTION_POST_JSON_END," + workflow_id + ",NA,1," + moment().format('MMMM Do YYYY, h:mm:ss a') + "\n");
    }
    res.send(reply)
  }
  else {
    var workflow_id = Math.floor(Math.random() * 10000000).toString();
    if ("LOGGING_NAME" in process.env) {
      client.append(process.env.LOGGING_NAME, "EXECUTION_POST_TEXT_START," + workflow_id + ",NA,1," + moment().format('MMMM Do YYYY, h:mm:ss a') + "\n");
    }
    var [reply, code] = func.function_handler({"request": req.body, "workflow_id": workflow_id, "request_type": "POST", "is_json": false});
    if ("LOGGING_NAME" in process.env) {
      client.append(process.env.LOGGING_NAME, "EXECUTION_POST_TEXT_END," + workflow_id + ",NA,1," + moment().format('MMMM Do YYYY, h:mm:ss a') + "\n");
    }
    res.send(reply)
  }
})

function gRPCFunctionHandler(call, callback) {
  if ("LOGGING_NAME" in process.env) {
    client.append(process.env.LOGGING_NAME, "EXECUTION_GRPC_START," + call.request.getWorkflow_id() + ",NA,1," + moment().format('MMMM Do YYYY, h:mm:ss a') + "\n");
  }
  var [reply, code] = func.function_handler({"request": call.request.getBody(), "workflow_id": call.request.getWorkflow_id(), "request_type": "GRPC", "is_json": false});
  if ("LOGGING_NAME" in process.env) {
    client.append(process.env.LOGGING_NAME, "EXECUTION_GRPC_END," + call.request.getWorkflow_id() + ",NA,1," + moment().format('MMMM Do YYYY, h:mm:ss a') + "\n");
  }
  var r = new messages.ResponseBody();
  r.setReply(reply);
  r.setCode(code);
  callback(null, r);
}

function main() {
  if (!("SERVICE_TYPE" in process.env) || process.env.SERVICE_TYPE == "HTTP") {
    app.listen(process.env.FUNC_PORT, () => {})
  }
  else if (process.env.SERVICE_TYPE == "GRPC") {
    var server = new grpc.Server();
    server.addService(services.gRPCFunctionService, {gRPCFunctionHandler: gRPCFunctionHandler});
    server.bindAsync('0.0.0.0:' + process.env.FUNC_PORT, grpc.ServerCredentials.createInsecure(), () => {
      server.start();
    });
  }
}

main();
