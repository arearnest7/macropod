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

// purpose: serve the function invoke and return the results to the appropriate request handler
// input:
//  request: MacroPodRequest object containing the payload and all invocation related data
// output:
//  MacroPodReply: reply string and code
async function Serve_Invoke(request) {
  var reply = "";
  var code = 500;
  await rpc.Timestamp(request, "", "request");
  [reply, code] = await func.FunctionHandler(request);
  await rpc.Timestamp(request, "", "request_end");
  var res = {Reply: reply, Code: code};
  return res;
}

// purpose: gRPC handler for Invoke
// input:
//  call: context object for request with payload
//  callback: return handler for gRPC nodejs
// output:
//  MacroPodReply: reply protobuf struct
async function Invoke(call, callback) {
  var j = {};
  if (call.request.JSON && call.request.JSON.fields) {
    for (var key in call.request.JSON.fields) {
      if (call.request.JSON.fields[key].stringValue !== undefined) {
        j[key] = call.request.JSON.fields[key].stringValue
      } else if (call.request.JSON.fields[key].numberValue !== undefined) {
        j[key] = call.request.JSON.fields[key].numberValue
      } else if (call.request.JSON.fields[key].boolValue !== undefined) {
        j[key] = call.request.JSON.fields[key].boolValue
      } else if (call.request.JSON.fields[key].nullValue !== undefined) {
        j[key] = call.request.JSON.fields[key].nullValue
      }
    }
  }
  var request = {
    Workflow: call.request.Workflow,
    Function: call.request.Function,
    WorkflowID: call.request.WorkflowID,
    Depth: call.request.Depth,
    Width: call.request.Width,
    Target: call.request.Target,
    JSON: j,
    Data: call.request.Data,
    Text: call.request.Text,
  };
  const res = await Serve_Invoke(request);
  await callback(null, res);
}

// purpose: http handler for Invoke
// input:
//  res: writer to provide HTTP reply data
//  req: GET/POST object
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

app.get('/', async (req, res) => { // bind GET requests to use HTTP_Invoke
  var reply = await HTTP_Invoke(req);
  res.send(reply["Reply"]);
})

app.post('/', async (req, res) => { // bind POST requests to use HTTP_Invoke
  var reply = await HTTP_Invoke(req);
  res.send(reply["Reply"]);
})

// input:
//  SERVICE_PORT: string defining gRPC service binding; default 5000
//  HTTP_PORT: string defining HTTP service binding; default 6000
//  WORKFLOW: string name of workflow the function belongs to
//  FUNCTION: string name of the function
//  INGRESS: string representing the IP:port of the ingress such that gateway communication can be used if enabled
//  LOGGER: string IP:port of the logger if logging is enabled; empty if logging is disabled
//  COMM_TYPE: string name of the communication type (refer to the deployer documentation)
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
