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

const RPC = (context, dest, payloads) => {
	if ("LOGGING_NAME" in process.env) {
		const client = redis.createClient({url: process.env.LOGGING_URL, password: process.env.LOGGING_PASSWORD});
        	client.append(process.env.LOGGING_NAME, moment().format('MMMM Do YYYY, h:mm:ss a') + "," + context.WorkflowId + "," + context.Depth.toString() + "," + context.Width.toString() + "," + context.RequestType + "," + "10" + "\n");
	}
	var stub = new gRPCFunction.gRPCFunction(dest, grpc.credentials.createInsecure());
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
	payloads.forEach(function(payload) {
		if (request_type == "gg" || request_type == "gm") {
			var request = {
				data: payload,
				workflow_id: context.WorkflowId,
				depth: context.Depth+1,
				width: i,
				request_type: request_type
			};
			stub.gRPCFunctionHandler(request, function(err, response) {
				console.log(response.reply);
				tl.push(response);
			});
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
                        stub.gRPCFunctionHandler(request, function(err, response) {
                                tl.push(response);
                        });
		}
		i += 1;
	});
	return [tl.length.toString()];
	/*var results = new Array();
	i = 0;
	if (!("RPC_DEST_PV" in process.env)) {
		tl.forEach(function(t) {
			results[i] = t.reply;
			i += 1;
		});
	} else {
		tl.forEach(function(t) {
			var reply = fs.readFileSync(process.env.RPC_DEST_PV + "/" + t.pv_path);
			fs.rmSync(process.env.RPC_DEST_PV + "/" + t.pv_path);
			results[i] = reply;
			i += 1;
		});
	}
	if ("RPC_PV" in process.env) {
		pv_paths.forEach(function(pv_path) {
			fs.rmSync(process.env.RPC_PV + "/" + pv_path);
		});
	}
	if ("LOGGING_NAME" in process.env) {
                client.append(process.env.LOGGING_NAME, moment().format('MMMM Do YYYY, h:mm:ss a') + "," + context.WorkflowId + "," + context.Depth.toString() + "," + context.Width.toString() + "," + context.RequestType + "," + "11" + "\n");
        }
	return results;*/
}

// Export the function
module.exports = { RPC };
