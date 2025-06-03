const grpc = require('@grpc/grpc-js');
const exec = require('child_process').execSync;
const moment = require('moment');
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
  await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + depth.toString() + "," + width.toString() + "," + request_type + "," + "request_start");
  ctx.Request = data;
  await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + depth.toString() + "," + width.toString() + "," + request_type + "," + "function_start");
  [reply_t, code_t] = await func.FunctionHandler(ctx);
  await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + depth.toString() + "," + width.toString() + "," + request_type + "," + "request_end");
  var res = {reply: reply_t, code: code_t, pv_path: pv_path_t};
  return res;
}

async function gRPCFunctionHandler(call, callback) {
  const res = await invoke(call.request);
  await callback(null, res);
}

function main() {
  var server = new grpc.Server({'grpc.max_send_message_length': 1024*1024*200, 'grpc.max_receive_message_length': 1024*1024*200});
  server.addService(gRPCFunction.gRPCFunction.service, {gRPCFunctionHandler: gRPCFunctionHandler});
  server.bindAsync('0.0.0.0:' + process.env.FUNC_PORT, grpc.ServerCredentials.createInsecure(), () => {});
}

main();
