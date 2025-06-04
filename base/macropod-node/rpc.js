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
  if process.env["COMM_TYPE"] == "direct" {
    return new Promise((resolve, reject) => {
      var stub = new macropod_pb.MacroPodFunction(process.env[dest], grpc.credentials.createInsecure());
      stub.Invoke(request, async function(err, response) {
        resolve(response);
      });
    });
  }
  else if process.env["COMM_TYPE"] == "gateway" {
    return new Promise((resolve, reject) => {
      var stub = new macropod_pb.MacroPodIngress(process.env[dest], grpc.credentials.createInsecure());
      stub.FunctionInvoke(request, async function(err, response) {
        resolve(response);
      });
    });
  }
  return new Promise((resolve, reject) => {
    var stub = new macropod_pb.MacroPodFunction(process.env[dest], grpc.credentials.createInsecure());
    stub.Invoke(request, async function(err, response) {
      resolve(response);
    });
  });
}

async function Invoke(context, dest, payload) {
  await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + context.WorkflowID + "," + context.Depth.toString() + "," + context.Width.toString() + ",invoke_rpc_start");
  var request = {
    Function: dest,
    Text: payload,
    WorkflowID: context.WorkflowID,
    Depth: context.Depth+1,
    Width: 0
  };
  let [reply, code] = await invoke(dest, request);
  await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + context.WorkflowID + "," + context.Depth.toString() + "," + context.Width.toString() + ",invoke_rpc_end");
  return [reply, code];
}

async function Invoke_JSON(context, dest, payload) {
  await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + context.WorkflowID + "," + context.Depth.toString() + "," + context.Width.toString() + ",invoke_rpc_start");
  var request = {
    Function: dest,
    JSON: payload,
    WorkflowID: context.WorkflowID,
    Depth: context.Depth+1,
    Width: 0
  };
  let [reply, code] = await invoke(dest, request);
  await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + context.WorkflowID + "," + context.Depth.toString() + "," + context.Width.toString() + ",invoke_rpc_end");
  return [reply, code];
}

async function Invoke_Data(context, dest, payload) {
  await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + context.WorkflowID + "," + context.Depth.toString() + "," + context.Width.toString() + ",invoke_rpc_start");
  var request = {
    Function: dest,
    Data: payload,
    WorkflowID: context.WorkflowID,
    Depth: context.Depth+1,
    Width: 0
  };
  let [reply, code] = await invoke(dest, request);
  await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + context.WorkflowID + "," + context.Depth.toString() + "," + context.Width.toString() + ",invoke_rpc_end");
  return [reply, code];
}

async function Invoke_Multi(context, dest, payloads) {
  await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + context.WorkflowID + "," + context.Depth.toString() + "," + context.Width.toString() + ",invoke_rpc_start");
  var tl = new Array();
  var i = 0;
  for (let payload of payloads) {
    var request = {
      Function: dest,
      Text: payload,
      WorkflowID: context.WorkflowID,
      Depth: context.Depth+1,
      Width: i
    };
    let t = await invoke(dest, request);
    tl.push(t);
    i += 1;
  }
  var reply = new Array();
  var code = new Array();
  for (let t of tl) {
    reply.push(t.Reply);
    code.push(t.Code)
  }
  await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + context.WorkflowID + "," + context.Depth.toString() + "," + context.Width.toString() + ",invoke_rpc_end");
  return [reply, code];
}

async function Invoke_Multi_JSON(context, dest, payloads) {
  await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + context.WorkflowID + "," + context.Depth.toString() + "," + context.Width.toString() + ",invoke_rpc_start");
  var tl = new Array();
  var i = 0;
  for (let payload of payloads) {
    var request = {
      Function: dest,
      JSON: payload,
      WorkflowID: context.WorkflowID,
      Depth: context.Depth+1,
      Width: i
    };
    let t = await invoke(dest, request);
    tl.push(t);
    i += 1;
  }
  var reply = new Array();
  var code = new Array();
  for (let t of tl) {
    reply.push(t.Reply);
    code.push(t.Code)
  }
  await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + context.WorkflowID + "," + context.Depth.toString() + "," + context.Width.toString() + ",invoke_rpc_end");
  return [reply, code];
}

async function Invoke_Multi_Data(context, dest, payloads) {
  await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + context.WorkflowID + "," + context.Depth.toString() +"," + context.Width.toString() + ",invoke_rpc_start");
  var tl = new Array();
  var i = 0;
  for (let payload of payloads) {
    var request = {
      Function: dest,
      Data: payload,
      WorkflowID: context.WorkflowID,
      Depth: context.Depth+1,
      Width: i
    };
    let t = await invoke(dest, request);
    tl.push(t);
    i += 1;
  }
  var reply = new Array();
  var code = new Array();
  for (let t of tl) {
    reply.push(t.Reply);
    code.push(t.Code)
  }
  await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + context.WorkflowID + "," + context.Depth.toString() + "," + context.Width.toString() + ",invoke_rpc_end");
  return [reply, code];
}

// Export the function
module.exports = { Invoke };
