const require('./rpc')
const redis = require('redis');

//const client = redis.createClient({url: process.env.REDIS_URL, password: process.env.REDIS_PASSWORD});

const function_handler = async (context) => {
	//client.on('error', err => console.log('Redis Client Error', err));
        //await client.connect();
	var body = JSON.parse(context["request"]);
	reply = 1; //await client.exists("voter-" + body['id']);
	if (reply == 1) {
		const g_val = "Not Voted"; //await client.get("voter-" + body['id']);
		if (g_val != "Not Voted") {
			return [body['id'] + " already submitted a vote.", 200];
		}
		else {
			data = '';
			newBody = body;
			rpc.RPC(process.env.ELECTION_VOTE_PROCESSOR, [newBody], context["workflow_id"])
			return ["Vote " + body['id'] + " registered", 200];
		}
	}
	return ["This voter id does not exist: " + body['id'], 200];
}

// Export the function
module.exports = { function_handler };

