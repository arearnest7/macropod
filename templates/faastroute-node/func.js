const rpc = require('./rpc')

async function FunctionHandler(context) {
	if (context.InvokeType != "GRPC") {
		const payloads = new Array();
        	payloads.push(Buffer.from("a".repeat(1000000), "utf8"));
		const res = await rpc.RPC(context, process.env.TEST, payloads);
		const res1 = await res.toString();
		return [res1, 200];
	}
	return [context.Request.toString(), 200];
}

// Export the function
module.exports = { FunctionHandler };
