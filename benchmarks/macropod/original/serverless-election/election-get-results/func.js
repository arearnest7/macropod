const rpc = require('./rpc')
const redis = require('redis');

//const client = redis.createClient({url: process.env.REDIS_URL, password: process.env.REDIS_PASSWORD});

const state_list = ['AK', 'AL', 'AR', 'AZ', 'CA', 'CO', 'CT', 'DC', 'DE', 'FL', 'GA', 'HI', 'IA', 'ID', 'IL', 'IN', 'KS', 'KY', 'LA', 'MA', 'MD', 'ME', 'MI', 'MN', 'MO', 'MS', 'MT', 'NC', 'ND', 'NE', 'NH', 'NJ', 'NM', 'NV', 'NY', 'OH', 'OK', 'OR', 'PA', 'RI', 'SC', 'SD', 'TN', 'TX', 'U'];

async function FunctionHandler(context) {
	//client.on('error', err => console.log('Redis Client Error', err));
        //await client.connect();
	var results = [];
	for (var state in state_list) {
		state_results = {}; //await client.keys('election-results-' + state + '-*');
		var total_count = {"total": 0};
		for (var i = 0; i < state_results.length; i++) {
			const cnt = 1; //await parseInt(client.get(state_results[i]));
			total_count[state_results[i]] = cnt;
			total_count["total"] += cnt;
		}

		var state_results_final = {};
		for (const [candidate, total] of total_count.entries()) {
			if (candidate !== "total") {
				state_results_final[candidate] = (total / total_count["total"] * 100.0).round();
			}
		}
		results.push({
			"state": state,
			"disclaimer": "These vote counts are estimates. Visit https://github.com/tylerpearson/serverless-election-aws for more info.",
			"total_count": total_count,
			"results": state_results_final
		});
	}
	//var response = {
	//	"isBase64Encoded": false,
    	//	"statusCode": 200,
    	//	"headers": {
      	//		"Access-Control-Allow-Origin": "*"
    	//	},
    	//	"body": results
	//};
	return [results.toString(), 200];
}

// Export the function
module.exports = { FunctionHandler };
