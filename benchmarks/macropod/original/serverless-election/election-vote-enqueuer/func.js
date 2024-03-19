const require('./rpc')
const redis = require('redis');

//const client = redis.createClient({url: process.env.REDIS_URL, password: process.env.REDIS_PASSWORD});

async function FunctionHandler(context) {
	//client.on('error', err => console.log('Redis Client Error', err));
        //await client.connect();
	var body = JSON.parse(context.Request);
	reply = 1; //await client.exists("voter-" + body['id']);
	if (reply == 1) {
		const g_val = "Not Voted"; //await client.get("voter-" + body['id']);
		if (g_val != "Not Voted") {
			return [body['id'] + " already submitted a vote.", 200];
		}
		else {
			data = '';
			newBody = body;
			rpc.RPC(context, process.env.ELECTION_VOTE_PROCESSOR, [context.Request]);
			return ["Vote " + body['id'] + " registered", 200];
		}
	}
	return ["This voter id does not exist: " + body['id'], 200];
}

// Export the function
module.exports = { FunctionHandler };

