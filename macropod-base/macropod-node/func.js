const rpc = require('./rpc')

async function FunctionHandler(context) {
	if (process.env.TEST) {
		const payloads = new Array();
        	payloads.push(Buffer.from("a".repeat(10000), "utf8"));
		const res = await rpc.RPC(context, process.env.TEST, payloads);
		const res1 = await res[0];
		return [res1, 200];
	}
	return [context.Request.toString(), 200];
}

// Export the function
module.exports = { FunctionHandler };
