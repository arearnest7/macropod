const rpc = require('./rpc')

async function FunctionHandler(context) {
	var body = context.JSON;
	if (body["requestType"] ==  'get_results') {
		let [res, code] = await rpc.Invoke_JSON(context, "ELECTION_GET_RESULTS", body);
		return [res, code];
	}
	else if (body["requestType"] == 'vote') {
		let [res, code] = await rpc.Invoke_JSON(context, "ELECTION_VOTE_ENQUEUER", body);
		return [res, code];
	}
	return ['invalid request type', 200];
}

// Export the function
module.exports = { FunctionHandler };
