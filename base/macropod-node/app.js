const express = require('express');
const app = express();
const func = require('./func');
const rpc = require('./rpc');
const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
var packageDefinition = protoLoader.loadSync(
  "./macropod.proto",
  { keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true
  });
var macropod_pb = grpc.loadPackageDefinition(packageDefinition).macropod;

async function Serve_Invoke(request) {
  var reply = "";
  var code = 500;
  await rpc.Timestamp(request, "", "request");
  [reply, code] = await func.FunctionHandler(request);
  await rpc.Timestamp(request, "", "request_end");
  var res = {reply: reply, code: code};
  return res;
}

async function Invoke(call, callback) {
  const res = await Serve_Invoke(call.request);
  await callback(null, res);
}

async function HTTP_Invoke(req) {
  var workflow_id = Math.floor(Math.random() * 100000).toString();
  var request = {
    Workflow: process.env.Workflow,
    Function: process.env.Function,
    WorkflowID: workflow_id,
    Depth: 0,
    Width: 0,
    Target: req.get("host"),
  };
  var content = req.get("Content-Type");
  var body = '';
  await req.on('data', chunk => {
    body += chunk.toString();
  });
  switch(content) {
    case "text/plain":
      request.Text = body;
      break;
    case "application/json":
      var j = JSON.parse(body);
      request.JSON = j;
      break;
    case "application/octet-stream":
      var d = Buffer.from(body);
      request.Data = d;
      break;
    default:
      request.Text = body;
      break;
  }
  const results = await Serve_Invoke(request);
  return results;
}

app.get('/', async (req, res) => {
  var reply = await HTTP_Invoke(req);
  res.send(reply);
})

app.post('/', async (req, res) => {
  var reply = await HTTP_Invoke(req);
  res.send(reply);
})

function main() {
  var service_port = process.env.SERVICE_PORT;
  if (service_port === "") {
    service_port = "5000"
  }
  var http_port = process.env.HTTP_PORT;
  if (http_port === "") {
    http_port = "6000"
  }
  var server = new grpc.Server({'grpc.max_send_message_length': 1024*1024*200, 'grpc.max_receive_message_length': 1024*1024*200});
  server.addService(macropod_pb.MacroPodFunction.service, {Invoke: Invoke});
  server.bindAsync('0.0.0.0:' + service_port, grpc.ServerCredentials.createInsecure(), () => {});
  app.listen(http_port, () => {})
}

main();
