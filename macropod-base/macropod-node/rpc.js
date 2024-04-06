const moment = require('moment')
const fs = require('fs')
const grpc = require('@grpc/grpc-js')
const protoLoader = require('@grpc/proto-loader');
var packageDefinition = protoLoader.loadSync(
  "./app.proto",
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
	console.log(moment().format('MMMM Do YYYY, h:mm:sss a') + "," + context.WorkflowId + "," + context.Depth.toString() + "," + context.Width.toString() + "," + context.RequestType + "," + "8" + "\n");
	var tl = new Array();
	var pv_paths = new Array();
	var request_type = "gg";
	if (!("RPC_PV" in process.env)) {
		if (!("RPC_DEST_PV" in process.env)) {
			request_type = "gg";
		}
		else {
			request_type = "gm";
		}
	}
	else {
		if (!("RPC_DEST_PV" in process.env)) {
			request_type = "mg";
		}
		else {
			request_type = "mm";
		}
	}
	var i = 0;
	for (let payload of payloads) {
		if (request_type == "gg" || request_type == "gm") {
			var request = {
				data: payload,
				workflow_id: context.WorkflowId,
				depth: context.Depth+1,
				width: i,
				request_type: request_type
			};
			let t = await invoke(dest, request);
			tl.push(t);
		} else {
			var pv_path = context.WorkflowId + "_" + context.Depth.toString() + "_" + context.Width.toString() + "_" + Math.floor(Math.random() * 10000000).toString();
			pv_paths.push(pv_path);
			fs.writeFileSync(process.env.RPC_PV + "/" + pv_path, payload);
			var request = {
                        	workflow_id: context.WorkflowId,
                        	depth: context.Depth+1,
				width: i,
                        	request_type: request_type,
				pv_path: pv_path
			};
                        let t = await invoke(dest, request);
                        tl.push(t);
		}
		i += 1;
	}
	var results = new Array();
	if (!("RPC_DEST_PV" in process.env)) {
		for (let t of tl) {
			console.log(t.reply);
			var reply = t.reply;
			results.push(reply);
		}
	} else {
		for (let t of tl) {
			var reply = fs.readFileSync(process.env.RPC_DEST_PV + "/" + t.pv_path);
			fs.rmSync(process.env.RPC_DEST_PV + "/" + t.pv_path);
			results.push(reply);
		}
	}
	if ("RPC_PV" in process.env) {
		for (let pv_path of pv_paths) {
			fs.rmSync(process.env.RPC_PV + "/" + pv_path);
		}
	}
        console.log(moment().format('MMMM Do YYYY, h:mm:sss a') + "," + context.WorkflowId + "," + context.Depth.toString() + "," + context.Width.toString() + "," + context.RequestType + "," + "9" + "\n");
	return results;
}

// Export the function
module.exports = { RPC };
