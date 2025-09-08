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
var macropod_pb = grpc.loadPackageDefinition(packageDefinition).macropod;

// purpose: send a message to logger (if enabled) to dictate a timestamp of function operation across multiple replicas.
// input:
//  context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
//  target: invocation target for the RPC (if it exists); otherwise null
//  message: message to be stored on the logger if logging is enabled
function Timestamp(context, target, message) {
  var request = {
    Workflow: context.Workflow,
    Function: context.Function,
    WorkflowID: context.WorkflowID,
    Depth: context.Depth,
    Width: context.Width,
    Target: target,
    Text: message
  };
  if (process.env.Logger && process.env.LOGGER != "") {
    var stub = new macropod_pb.MacroPodLogger(process.env.LOGGER, grpc.credentials.createInsecure());
    stub.Timestamp(request, async function(err, response) {});
  }
}

// purpose: send a error to logger (if enabled) notifying that a replica reports an error of err. this is to prevent critical errors for developers/debuggers from being lost due to replica deletion or restart.
// input:
//  context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
//  err: error message to be stored on the logger if logging is enabled
function Error(context, err) {
  var request = {
    Workflow: context.Workflow,
    Function: context.Function,
    Text: err
  };
  if (process.env.Logger && process.env.LOGGER != "") {
    var stub = new macropod_pb.MacroPodLogger(process.env.LOGGER, grpc.credentials.createInsecure());
    stub.Error(request, async function(err, response) {});
  }
}

// purpose: send a message to logger (if enabled) to allow persistent printouts for developers/debuggers.
// input:
//  context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
//  message: message to be stored on the logger if logging is enabled
function Print(context, message) {
  var request = {
    Workflow: context.Workflow,
    Function: context.Function,
    Text: message
  };
  if (process.env.Logger && process.env.LOGGER != "") {
    var stub = new macropod_pb.MacroPodLogger(process.env.LOGGER, grpc.credentials.createInsecure());
    stub.Print(request, async function(err, response) {});
  }
}

function invoke(dest, request) { // general purpose RPC handler used by the other async functions to provide the specific RPC behavior desired
  if (process.env.COMM_TYPE == "direct") {  // invocation sent directly to target bypassing gateway (direct RPC or IPC)
    return new Promise((resolve, reject) => {
      var stub = new macropod_pb.MacroPodFunction(process.env[dest], grpc.credentials.createInsecure());
      stub.Invoke(request, async function(err, response) {
        await Error(request, err);
        resolve(response);
      });
    });
  } else if (process.env.COMM_TYPE == "gateway") { // invocation sent to relay off of gateway (standard serverless behavior)
    return new Promise((resolve, reject) => {
      var stub = new macropod_pb.MacroPodIngress(process.env.INGRESS, grpc.credentials.createInsecure());
      stub.FunctionInvoke(request, async function(err, response) {
        await Error(request, err);
        resolve(response);
      });
    });
  }
  return new Promise((resolve, reject) => {
    var stub = new macropod_pb.MacroPodFunction(process.env[dest], grpc.credentials.createInsecure());
    stub.Invoke(request, async function(err, response) {
      await Error(request, err);
      resolve(response);
    });
  });
}

// purpose: invoke a downstream function with a text/string payload and return its string output
// input:
//  context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
//  dest: invocation target for the RPC; this is the name of function 
//  payload: string payload to be sent to destination function
// output: 
//  string: RPC return output
//  int: RPC error code
async function Invoke(context, dest, payload) {
  await Timestamp(context, dest, "Invoke");
  var request = {
    Workflow: context.Workflow,
    Function: dest,
    Target: process.env[dest],
    Text: payload,
    WorkflowID: context.WorkflowID,
    Depth: context.Depth+1,
    Width: 0
  };
  let response = await invoke(dest, request);
  await Timestamp(context, dest, "Invoke_End");
  return [response.Reply, response.Code];
}

// purpose: invoke a downstream function with a JSON list (dictionary) payload and return its string output
// input:
//  context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
//  dest: invocation target for the RPC; this is the name of function 
//  payload: dictionary payload to be sent to destination function
// output: 
//  string: RPC return output
//  int: RPC error code
async function Invoke_JSON(context, dest, payload) {
  await Timestamp(context, dest, "Invoke_JSON");
  payload_s = {fields: {}};
  for (var key in payload) {
    payload_s.fields[key] = {};
    if (typeof payload[key] === "string") {
      payload_s.fields[key].stringValue = payload[key];
    } else if (typeof payload[key] === "number") {
      payload_s.fields[key].numberValue = payload[key];
    } else if (typeof payload[key] === "bool") {
      payload_s.fields[key].boolValue = payload[key];
    } else if (typeof payload[key] === "null") {
      payload_s.fields[key].nullValue = payload[key];
    }
  }
  var request = {
    Workflow: context.Workflow,
    Function: dest,
    Target: process.env[dest],
    JSON: payload_s,
    WorkflowID: context.WorkflowID,
    Depth: context.Depth+1,
    Width: 0
  };
  let response = await invoke(dest, request);
  await Timestamp(context, dest, "Invoke_JSON_End");
  return [response.Reply, response.Code];
}

// purpose: invoke a downstream function with a bytes payload and return its string output
// input:
//  context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
//  dest: invocation target for the RPC; this is the name of function 
//  payload: bytes payload to be sent to destination function
// output: 
//  string: RPC return output
//  int: RPC error code
async function Invoke_Data(context, dest, payload) {
  await Timestamp(context, dest, "Invoke_Data");
  var request = {
    Workflow: context.Workflow,
    Function: dest,
    Target: process.env[dest],
    Data: payload,
    WorkflowID: context.WorkflowID,
    Depth: context.Depth+1,
    Width: 0
  };
  let response = await invoke(dest, request);
  await Timestamp(context, dest, "Invoke_Data_End");
  return [response.Reply, response.Code];
}

// purpose: invoke a downstream function multiple times concurrently with a list of text/string payloads and return all string outputs
// input:
//  context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
//  dest: invocation target for the RPCs; this is the name of function 
//  payloads: list of string payloads to be sent to destination function
// output: 
//  string array: array of RPC return outputs
//  int array: array of RPC error code
async function Invoke_Multi(context, dest, payloads) {
  await Timestamp(context, dest, "Invoke_Mutli");
  var tl = new Array();
  var i = 0;
  for (let payload of payloads) {
    var request = {
      Workflow: context.Workflow,
      Function: dest,
      Target: process.env[dest],
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
  await Timestamp(context, dest, "Invoke_Multi_End");
  return [reply, code];
}

// purpose: invoke a downstream function multiple times concurrently with a list of JSON list (dictionary) payloads and return all string outputs
// input:
//  context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
//  dest: invocation target for the RPCs; this is the name of function 
//  payloads: list of dictionary payloads to be sent to destination function
// output: 
//  string array: array of RPC return outputs
//  int array: array of RPC error code
async function Invoke_Multi_JSON(context, dest, payloads) {
  await Timestamp(context, dest, "Invoke_Multi_JSON");
  var tl = new Array();
  var i = 0;
  for (let payload of payloads) {
    payload_s = {fields: {}};
    for (var key in payload) {
      payload_s.fields[key] = {};
      if (typeof payload[key] === "string") {
        payload_s.fields[key].stringValue = payload[key];
      } else if (typeof payload[key] === "number") {
        payload_s.fields[key].numberValue = payload[key];
      } else if (typeof payload[key] === "bool") {
        payload_s.fields[key].boolValue = payload[key];
      } else if (typeof payload[key] === "null") {
        payload_s.fields[key].nullValue = payload[key];
      }
    }
    var request = {
      Workflow: context.Workflow,
      Function: dest,
      Target: process.env[dest],
      JSON: payload_s,
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
  await Timestamp(context, dest, "Invoke_Multi_JSON_End");
  return [reply, code];
}

// purpose: invoke a downstream function multiple times concurrently with a list of bytes payloads and return all string outputs
// input:
//  context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
//  dest: invocation target for the RPCs; this is the name of function 
//  payloads: list of bytes payloads to be sent to destination function
// output: 
//  string array: array of RPC return outputs
//  int array: array of RPC error code
async function Invoke_Multi_Data(context, dest, payloads) {
  await Timestamp(context, dest, "Invoke_Multi_Data");
  var tl = new Array();
  var i = 0;
  for (let payload of payloads) {
    var request = {
      Workflow: context.Workflow,
      Function: dest,
      Target: process.env[dest],
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
  await Timestamp(context, dest, "Invoke_Multi_Data_End");
  return [reply, code];
}

// Export the function
module.exports = { Timestamp, Error, Print, Invoke, Invoke_JSON, Invoke_Data, Invoke_Multi, Invoke_Multi_JSON, Invoke_Multi_Data };
