const grpc = require('@grpc/grpc-js');
const exec = require('child_process').execSync;
const moment = require('moment');
const func = require('./func');
const protoLoader = require('@grpc/proto-loader');
var packageDefinition = protoLoader.loadSync(
  "./macropod.proto",
  { keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true
  });
var macropod_pb = grpc.loadPackageDefinition(packageDefinition).function;

async function invoke(request) {
  var reply = "";
  var code = 500;
  await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + request.WorkflowID + "," + request.Depth.toString() + "," + request.Width.toString() + ",request_start");
  [reply, code] = await func.FunctionHandler(request);
  await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + request.WorkflowID + "," + request.Depth.toString() + "," + request.Width.toString() + ",request_end");
  var res = {reply: reply, code: code};
  return res;
}

async function FunctionInvoke(call, callback) {
  const res = await invoke(call.request);
  await callback(null, res);
}

function main() {
  var server = new grpc.Server({'grpc.max_send_message_length': 1024*1024*200, 'grpc.max_receive_message_length': 1024*1024*200});
  server.addService(macropod_pb.MacroPodFunction.service, {Invoke: FunctionInvoke});
  server.bindAsync('0.0.0.0:' + process.env.FUNC_PORT, grpc.ServerCredentials.createInsecure(), () => {});
}

main();
