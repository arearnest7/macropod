const rpc = require('./rpc')

const FunctionHandler = (context) => {
	const payloads = new Array();
	payloads.push("a".repeat(100));
	if (context.RequestType != "GRPC") {
		return ["[" + rpc.RPC(context, process.env.TEST, payloads).toString() + "]", 200];
	}
	return [context.Request.toString(), 200];
}

// Export the function
module.exports = { FunctionHandler };
