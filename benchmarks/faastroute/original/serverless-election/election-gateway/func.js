const require('./rpc')

const function_handler = async (context) => {
	if (context["request_type"] != "GRPC") {
		var body = context["request"];
		if (body['requestType'] ==  'get_results') {
			return [rpc.RPC(context, process.env.ELECTION_GET_RESULTS, [body])[0], 200];
		}
		else if (body['requestType'] == 'vote') {
			return [rpc.RPC(context, process.env.ELECTION_VOTE_ENQUEUER, [body])[0], 200];
		}
		return ['invalid request type', 200];
	}
}

// Export the function
module.exports = { function_handler };
