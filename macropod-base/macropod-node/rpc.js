const exec = require('child_process').execSync;
const moment = require('moment')
const grpc = require('@grpc/grpc-js')
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

function invoke(dest, request) {
	return new Promise((resolve, reject) => {
		var stub = new gRPCFunction.gRPCFunction(dest, grpc.credentials.createInsecure());
		stub.gRPCFunctionHandler(request, async function(err, response) {
			resolve(response);
		});
	});
}

async function RPC(context, dest, payloads) {
	await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + context.WorkflowId + "," + context.Depth.toString() + "," + context.Width.toString() + "," + context.RequestType + "," + "rpc_start");
	var tl = new Array();
	var pv_paths = new Array();
	var request_type = "gg";
	var i = 0;
	for (let payload of payloads) {
		var request = {
			data: payload,
			workflow_id: context.WorkflowId,
			depth: context.Depth+1,
			width: i,
			request_type: request_type
		};
		let t = await invoke(dest, request);
		tl.push(t);
		i += 1;
	}
	var results = new Array();
	for (let t of tl) {
		var reply = t.reply;
		results.push(reply);
	}
        await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + context.WorkflowId + "," + context.Depth.toString() + "," + context.Width.toString() + "," + context.RequestType + "," + "rpc_end");
	return results;
}

// Export the function
module.exports = { RPC };
