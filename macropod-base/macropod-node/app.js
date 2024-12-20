const express = require('express');
const app = express();
const grpc = require('@grpc/grpc-js');
const exec = require('child_process').execSync;
const moment = require('moment');
const fs = require('fs');
const func = require('./func');
const protoLoader = require('@grpc/proto-loader');
var packageDefinition = protoLoader.loadSync(
  "./wf.proto",
  { keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true
  });
var gRPCFunction = grpc.loadPackageDefinition(packageDefinition).function;
app.use(express.json());

app.get('/', async (req, res) => {
  var workflow_id = Math.floor(Math.random() * 10000000).toString();
  await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + "0" + "," + "0" + "," + "HTTP" + "," + "0");
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
    InvokeType: "HTTP",
    IsJson: true
  };
  if (req.header("Content-Type") == "application/json") {
    ctx.WorkflowId = workflow_id;
    ctx.IsJson = true;
    await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + "0" + "," + "0" + "," + "HTTP" + "," + "1");
    [reply, code] = await func.FunctionHandler(ctx);
    await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + "0" + "," + "0" + "," + "HTTP" + "," + "2");
    res.send(reply);
  }
  else {
    ctx.WorkflowId = workflow_id;
    ctx.IsJson = false;
    await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + "0" + "," + "0" + "," + "HTTP" + "," + "1");
    [reply, code] = await func.FunctionHandler(ctx);
    await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + "0" + "," + "0" + "," + "HTTP" + "," + "2");
    res.send(reply);
  }
})

app.post('/', async (req, res) => {
  var workflow_id = Math.floor(Math.random() * 10000000).toString();
  await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + "0" + "," + "0" + "," + "HTTP" + "," + "0");
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
    InvokeType: "HTTP",
    IsJson: true
  };
  if (req.header("Content-Type") == "application/json") {
    ctx.WorkflowId = workflow_id;
    ctx.IsJson = true;
    await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + "0" + "," + "0" + "," + "HTTP" + "," + "1");
    [reply, code] = await func.FunctionHandler(ctx);
    await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + "0" + "," + "0" + "," + "HTTP" + "," + "2");
    res.send(reply);
  }
  else {
    ctx.WorkflowId = workflow_id;
    ctx.IsJson = false;
    await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + "0" + "," + "0" + "," + "HTTP" + "," + "1");
    [reply, code] = await func.FunctionHandler(ctx);
    await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + "0" + "," + "0" + "," + "HTTP" + "," + "2");
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
  await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + depth.toString() + "," + width.toString() + "," + request_type + "," + "0");
  if (request_type == "" || request_type == "gg") {
    ctx.Request = data;
    await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + depth.toString() + "," + width.toString() + "," + request_type + "," + "1");
    [reply_t, code_t] = await func.FunctionHandler(ctx);
  }
  else if (request_type == "mg") {
    var req = fs.readFileSync(process.env.APP_PV + "/" + path);
    ctx.Request = req;
    await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + depth.toString() + "," + width.toString() + "," + request_type + "," + "1");
    [reply_t, code_t] = await func.FunctionHandler(ctx);
  }
  else if (request_type == "gm") {
    ctx.Request = data;
    var payload = "";
    await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + depth.toString() + "," + width.toString() + "," + request_type + "," + "1");
    [payload, code_t] = await func.FunctionHandler(ctx);
    pv_path_t = workflow_id + "_" + depth.toString() + "_" + width.toString() + "_" + Math.floor(Math.random() * 10000000).toString();
    fs.writeFileSync(process.env.APP_PV + "/" + pv_path_t, payload);
  }
  else if (request_type == "mm") {
    var req = fs.readFileSync(process.env.APP_PV + "/" + path);
    ctx.Request = req;
    var payload = "";
    await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + depth.toString() + "," + width.toString() + "," + request_type + "," + "1");
    [payload, code_t] = await func.FunctionHandler(ctx);
    pv_path_t = workflow_id + "_" + depth.toString() + "_" + width.toString() + "_" + Math.floor(Math.random() * 10000000).toString();
    fs.writeFileSync(process.env.APP_PV + "/" + pv_path_t, payload);
  }
  await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + depth.toString() + "," + width.toString() + "," + request_type + "," + "2");
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
