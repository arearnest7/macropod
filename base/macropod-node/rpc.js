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
  if process.env.LOGGER != "" {
    var stub = new macropod_pb.MacroPodLogger(process.env.LOGGER, grpc.credentials.createInsecure());
    stub.Timestamp(request, async function(err, response) {});
  }
}

function Error(context, err) {
  var request = {
    Workflow: context.Workflow,
    Function: context.Function,
    Text: err
  };
  if process.env.LOGGER != "" {
    var stub = new macropod_pb.MacroPodLogger(process.env.LOGGER, grpc.credentials.createInsecure());
    stub.Error(request, async function(err, response) {});
  }
}

function Print(context, message) {
  var request = {
    Workflow: context.Workflow,
    Function: context.Function,
    Text: message
  };
  if process.env.LOGGER != "" {
    var stub = new macropod_pb.MacroPodLogger(process.env.LOGGER, grpc.credentials.createInsecure());
    stub.Print(request, async function(err, response) {});
  }
}

function invoke(dest, request) {
  if process.env.COMM_TYPE == "direct" {
    return new Promise((resolve, reject) => {
      var stub = new macropod_pb.MacroPodFunction(process.env[dest], grpc.credentials.createInsecure());
      stub.Invoke(request, async function(err, response) {
        await Error(request, err);
        resolve(response);
      });
    });
  }
  else if process.env.COMM_TYPE == "gateway" {
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
  let [reply, code] = await invoke(dest, request);
  await Timestamp(context, dest, "Invoke_End");
  return [reply, code];
}

async function Invoke_JSON(context, dest, payload) {
  await Timestamp(context, dest, "Invoke_JSON");
  var request = {
    Workflow: context.Workflow,
    Function: dest,
    Target: process.env[dest],
    JSON: payload,
    WorkflowID: context.WorkflowID,
    Depth: context.Depth+1,
    Width: 0
  };
  let [reply, code] = await invoke(dest, request);
  await Timestamp(context, dest, "Invoke_JSON_End");
  return [reply, code];
}

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
  let [reply, code] = await invoke(dest, request);
  await Timestamp(context, dest, "Invoke_Data_End");
  return [reply, code];
}

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

async function Invoke_Multi_JSON(context, dest, payloads) {
  await Timestamp(context, dest, "Invoke_Multi_JSON");
  var tl = new Array();
  var i = 0;
  for (let payload of payloads) {
    var request = {
      Workflow: context.Workflow,
      Function: dest,
      Target: process.env[dest],
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
  await Timestamp(context, dest, "Invoke_Multi_JSON_End");
  return [reply, code];
}

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
module.exports = { Invoke };
