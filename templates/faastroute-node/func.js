const require('./rpc')

const function_handler = async (context) => {
	if (context["request_type"] != "GRPC"):
		return rpc.RPC(process.env.TEST, ["TEST"], context["workflow_id"]).toString(), 200
	return context["request"], 200
}

// Export the function
module.exports = { function_handler };
