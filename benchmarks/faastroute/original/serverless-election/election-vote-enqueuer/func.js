const require('./rpc')
const redis = require('redis');
const axios = require('axios');

//const client = redis.createClient({url: process.env.REDIS_URL, password: process.env.REDIS_PASSWORD});

const function_handler = async (context) => {
	//client.on('error', err => console.log('Redis Client Error', err));
        //await client.connect();
	reply = 1; //await client.exists("voter-" + body['id']);
	if (reply == 1) {
		const g_val = "Not Voted"; //await client.get("voter-" + body['id']);
		if (g_val != "Not Voted") {
			return body['id'] + " already submitted a vote.";
		}
		else {
			data = '';
			newBody = body;
			await axios.post(process.env.ELECTION_VOTE_PROCESSOR, newBody)
				.then( (response) => {
                                	data = response.data;
				});
			return "Vote " + body['id'] + " registered";
		}
	}
	return "This voter id does not exist: " + body['id'];
}

// Export the function
module.exports = { function_handler };

