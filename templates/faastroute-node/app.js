const express = require('express');
const app = express();
const grpc = require('@grpc/grpc-js');
const redis = require('redis');
const moment = require('moment');
const fs = require('fs');
const func = require('./func');
const protoLoader = require('@grpc/proto-loader');
var packageDefinition = protoLoader.loadSync(
  "./app.proto",
  { keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true
  });
var gRPCFunction = grpc.loadPackageDefinition(packageDefinition).function;

if ("LOGGING_NAME" in process.env) {
  const client = redis.createClient({url: process.env.LOGGING_URL, password: process.env.LOGGING_PASSWORD});
}

app.get('/', async (req, res) => {
  var request_type = "gg";
  if ("APP_PV" in process.env) {
    request_type = "gm";
  }
  var reply = "";
  var code = 500;
  var ctx = {
    Request: req.body,
    WorkflowId: "",
    Depth: 0,
    Width: 0,
    RequestType: request_type,
    InvokeType: "GET",
    IsJson: true
  };
  if (req.header("Content-Type") == "application/json") {
    var workflow_id = Math.floor(Math.random() * 10000000).toString();
    if ("LOGGING_NAME" in process.env) {
      await client.append(process.env.LOGGING_NAME, moment().format('MMMM Do YYYY, h:mm:ss a') + "," + workflow_id + "," + "0" + "," + "0" + "," + request_type + "," + "0" + "\n");
    }
    ctx.WorkflowId = workflow_id;
    ctx.IsJson = true;
    [reply, code] = await func.FunctionHandler(ctx);
    if ("LOGGING_NAME" in process.env) {
      await client.append(process.env.LOGGING_NAME, moment().format('MMMM Do YYYY, h:mm:ss a') + "," + workflow_id + "," + "0" + "," + "0" + "," + request_type + "," + "1" + "\n");
    }
    res.send(reply);
  }
  else {
    var workflow_id = Math.floor(Math.random() * 10000000).toString();
    if ("LOGGING_NAME" in process.env) {
      await client.append(process.env.LOGGING_NAME, moment().format('MMMM Do YYYY, h:mm:ss a') + "," + workflow_id + "," + "0" + "," + "0" + "," + request_type + "," + "2" + "\n");
    }
    ctx.WorkflowId = workflow_id;
    ctx.IsJson = false;
    [reply, code] = await func.FunctionHandler(ctx);
    if ("LOGGING_NAME" in process.env) {
      await client.append(process.env.LOGGING_NAME, moment().format('MMMM Do YYYY, h:mm:ss a') + "," + workflow_id + "," + "0" + "," + "0" + "," + request_type + "," + "3" + "\n");
    }
    res.send(reply);
  }
})

app.post('/', async (req, res) => {
  var request_type = "gg";
  if ("APP_PV" in process.env) {
    request_type = "gm";
  }
  var reply = "";
  var code = 500;
  var ctx = {
    Request: req.body,
    WorkflowId: "",
    Depth: 0,
    Width: 0,
    RequestType: request_type,
    InvokeType: "POST",
    IsJson: true
  };
  if (req.header("Content-Type") == "application/json") {
    var workflow_id = Math.floor(Math.random() * 10000000).toString();
    if ("LOGGING_NAME" in process.env) {
      await client.append(process.env.LOGGING_NAME, moment().format('MMMM Do YYYY, h:mm:ss a') + "," + workflow_id + "," + "0" + "," + "0" + "," + request_type + "," + "4" + "\n");
    }
    ctx.WorkflowId = workflow_id;
    ctx.IsJson = true;
    [reply, code] = await func.FunctionHandler(ctx);
    if ("LOGGING_NAME" in process.env) {
      await client.append(process.env.LOGGING_NAME, moment().format('MMMM Do YYYY, h:mm:ss a') + "," + workflow_id + "," + "0" + "," + "0" + "," + request_type + "," + "5" + "\n");
    }
    res.send(reply);
  }
  else {
    var workflow_id = Math.floor(Math.random() * 10000000).toString();
    if ("LOGGING_NAME" in process.env) {
      await client.append(process.env.LOGGING_NAME, moment().format('MMMM Do YYYY, h:mm:ss a') + "," + workflow_id + "," + "0" + "," + "0" + "," + request_type + "," + "6" + "\n");
    }
    ctx.WorkflowId = workflow_id;
    ctx.IsJson = false;
    [reply, code] = await func.FunctionHandler(ctx);
    if ("LOGGING_NAME" in process.env) {
      await client.append(process.env.LOGGING_NAME, moment().format('MMMM Do YYYY, h:mm:ss a') + "," + workflow_id + "," + "0" + "," + "0" + "," + request_type + "," + "7" + "\n");
    }
    res.send(reply);
  }
})

async function invoke(request) {
  var reply_t = "";
  var code_t = 500;
  var pv_path_t = "";
  var data = request.data;
  var workflow_id = request.workflow_id;
  var depth = request.depth;
  var width = request.width;
  var request_type = request.request_type;
  var path = request.pv_path;
  var ctx = {
    Request: "",
    WorkflowId: workflow_id,
    Depth: depth,
    Width: width,
    RequestType: request_type,
    InvokeType: "GRPC",
    IsJson: false
  };
  if ("LOGGING_NAME" in process.env) {
    await client.append(process.env.LOGGING_NAME, moment().format('MMMM Do YYYY, h:mm:ss a') + "," + workflow_id + "," + depth.toString() + "," + width.toString() + "," + request_type + "," + "8" + "\n");
  }
  if (request_type == "" || request_type == "gg") {
    ctx.Request = data;
    [reply_t, code_t] = await func.FunctionHandler(ctx);
  }
  else if (request_type == "mg") {
    var req = fs.readFileSync(process.env.APP_PV + "/" + path);
    ctx.Request = req;
    [reply_t, code_t] = await func.FunctionHandler(ctx);
  }
  else if (request_type == "gm") {
    ctx.Request = data;
    var payload = "";
    [payload, code_t] = await func.FunctionHandler(ctx);
    pv_path_t = workflow_id + "_" + depth.toString() + "_" + width.toString() + "_" + Math.floor(Math.random() * 10000000).toString();
    fs.writeFileSync(process.env.APP_PV + "/" + pv_path_t, payload);
  }
  else if (request_type == "mm") {
    var req = fs.readFileSync(process.env.APP_PV + "/" + path);
    ctx.Request = req;
    var payload = "";
    [payload, code_t] = await func.FunctionHandler(ctx);
    pv_path_t = workflow_id + "_" + depth.toString() + "_" + width.toString() + "_" + Math.floor(Math.random() * 10000000).toString();
    fs.writeFileSync(process.env.APP_PV + "/" + pv_path_t, payload);
  }
  if ("LOGGING_NAME" in process.env) {
    await client.append(process.env.LOGGING_NAME, moment().format('MMMM Do YYYY, h:mm:ss a') + "," + workflow_id + "," + depth.toString() + "," + width.toString() + "," + request_type + "," + "9" + "\n");
  }
  var res = {reply: reply_t, code: code_t, pv_path: pv_path_t};
  return res;
}

async function gRPCFunctionHandler(call, callback) {
  const res = await invoke(call.request);
  await callback(null, res);
}

function main() {
  if (!("SERVICE_TYPE" in process.env) || process.env.SERVICE_TYPE == "HTTP") {
    app.listen(process.env.FUNC_PORT, () => {});
  }
  else if (process.env.SERVICE_TYPE == "GRPC") {
    var server = new grpc.Server({'grpc.max_send_message_length': 1024*1024*200, 'grpc.max_receive_message_length': 1024*1024*200});
    server.addService(gRPCFunction.gRPCFunction.service, {gRPCFunctionHandler: gRPCFunctionHandler});
    server.bindAsync('0.0.0.0:' + process.env.FUNC_PORT, grpc.ServerCredentials.createInsecure(), () => {
      console.log("GRPC started");
    });
  }
}

main();
