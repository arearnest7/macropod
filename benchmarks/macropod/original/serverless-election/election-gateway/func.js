const rpc = require('./rpc')

async function FunctionHandler(context) {
	var body = JSON.parse(context.Request);
	if (body.requestType ==  'get_results') {
		var res = await rpc.RPC(context, process.env.ELECTION_GET_RESULTS, [context.Request]);
		return [res[0].toString(), 200];
	}
	else if (body.requestType == 'vote') {
		var res = await rpc.RPC(context, process.env.ELECTION_VOTE_ENQUEUER, [context.Request]);
		return [res[0].toString(), 200];
	}
	return ['invalid request type', 200];
}

// Export the function
module.exports = { FunctionHandler };
