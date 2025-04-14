const exec = require('child_process').execSync;
const moment = require('moment')
const grpc = require('@grpc/grpc-js')
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

function invoke(dest, request) {
  return new Promise((resolve, reject) => {
    var stub = new macropod_pb.MacroPodFunction(process.env[dest], grpc.credentials.createInsecure());
    stub.Invoke(request, async function(err, response) {
      resolve(response);
    });
  });
}

async function Invoke(context, dest, byte_payloads, string_payloads) {
  await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + context.InvocationID + "," + context.Depth.toString() + "," + context.Width.toString() + "," + "invoke_rpc_start");
  var tl = new Array();
  var i = 0;
  for (let payload of byte_payloads) {
    var request = {
      Function: dest,
      Data: payload,
      InvocationID: context.InvocationID,
      Depth: context.Depth+1,
      Width: i
    };
    let t = await invoke(dest, request);
    tl.push(t);
    i += 1;
  }
  for (let payload of string_payloads) {
    var request = {
      Function: dest,
      Text: payload,
      InvocationID: context.InvocationID,
      Depth: context.Depth+1,
      Width: i
    };
    let t = await invoke(dest, request);
    tl.push(t);
    i += 1;
  }
  var replies = new Array();
  var codes = new Array();
  for (let t of tl) {
    var reply = t.Reply;
    var code = t.Code;
    replies.push(reply);
    codes.push(code)
  }
  await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + context.InvocationID + "," + context.Depth.toString() + "," + context.Width.toString() + "," + "invoke_rpc_end");
  return [results, codes];
}

// Export the function
module.exports = { Invoke };
