const moment = require('moment')
const messages = require('./app_pb')
const services = require('./app_grpc_pb')
const grpc = require('@grpc/grpc-js')

const RPC = async (dest, payloads, id) => {
	if ("LOGGING_NAME" in process.env) {
		const client = redis.createClient({url: process.env.LOGGING_URL, password: process.env.LOGGING_PASSWORD});
        	client.append(process.env.LOGGING_NAME, "INVOKE_START," + id + "," + dest + "," + payloads.length.toString() + "," + moment().format('MMMM Do YYYY, h:mm:ss a') + "\n");
	}
	var stub = new services.gRPCFunctionClient(dest, grpc.credentials.createInsecure());
	var tl = [];
	payloads.forEach(function(payloads) {
		var request = new messages.RequestBody();
		request.setBody(payload);
		request.setWorkflow_id(id);
		stub.gRPCFunctionHandler(request, function(err, response) {
			tl.push(response);
		};
	}
        if ("LOGGING_NAME" in process.env) {
		client.append(process.env.LOGGING_NAME, "INVOKE_END," + id + "," + dest + "," + payloads.length.toString() + "," + moment().format('MMMM Do YYYY, h:mm:ss a') + "\n");
	}
	var results = [];
	tl.forEach(function(t) {
		results.push(t.getReply());
	}
	return results;
}

// Export the function
module.exports = { RPC };
