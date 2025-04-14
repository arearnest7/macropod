const rpc = require('./rpc')
const redis = require('redis');
const axios = require('axios');
const exec = require('child_process').execSync;
const moment = require('moment');

//const client = redis.createClient({url: process.env.REDIS_URL, password: process.env.REDIS_PASSWORD});

const state_list = ['AK', 'AL', 'AR', 'AZ', 'CA', 'CO', 'CT', 'DC', 'DE', 'FL', 'GA', 'HI', 'IA', 'ID'
, 'IL', 'IN', 'KS', 'KY', 'LA', 'MA', 'MD', 'ME', 'MI', 'MN', 'MO', 'MS', 'MT', 'NC', 'ND', 'NE', 'NH'
, 'NJ', 'NM', 'NV', 'NY', 'OH', 'OK', 'OR', 'PA', 'RI', 'SC', 'SD', 'TN', 'TX', 'U'];

const vote_processor_handler = async (body) => {
        //await client.set("voter-" + body['id'], JSON.stringify(body));

        var state = body['state'];
        var candidate = body['candidate'];

	reply = 1; //await client.exists("election-results-" + state + "-" + candidate);
        if (reply == 1) {
        	var cnt = 1; //parseInt(await client.get("election-results-" + state + "-" + candidate));
		cnt = cnt + 1;
		//await client.set("election-results-" + state + "-" + candidate, cnt.toString());
	}
	//else {
	//	await client.set("election-results-" + state + "-" + candidate, "1");
	//}
        return "success";
}

const vote_enqueuer_handler = async (body) => {
	reply = 1; //await client.exists("voter-" + body['id']);
	if (reply == 1) {
        	const g_val = "Not Voted" //await client.get("voter-" + body['id']);
		if (g_val != "Not Voted") {
			return {"isBase64Encoded": false, "statusCode": 409, "body": {"success": false, "message": (body['id'] + " already submitted a vote.")}};
		} else {
			let data = await vote_processor_handler(body);
			return {"isBase64Encoded": false, "statusCode": 201, "body": {"success": true, "message": ("Vote " + body['id'] + " registered")}};
		}
	}
	return {"isBase64Encoded": false, "statusCode": 404, "body": {"success": false, "message": ("This voter id does not exist: " + body['id'])}};
}

const get_results_handler = async (body) => {
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
        var response = {
                "isBase64Encoded": false,
                "statusCode": 200,
                "headers": {
                        "Access-Control-Allow-Origin": "*"
                },
                "body": results
        };
        return response;
}

async function FunctionHandler(context) {
	var body = JSON.parse(context.Request);
	if (body.requestType ==  'get_results') {
		//var res = await rpc.RPC(context, process.env.ELECTION_GET_RESULTS, [context.Request]);
		//return [res[0].toString(), 200];
		let data = await get_results_handler(newbody);
		return [data, 200];
	}
	else if (body.requestType == 'vote') {
		//var res = await rpc.RPC(context, process.env.ELECTION_VOTE_ENQUEUER, [context.Request]);
		//return [res[0].toString(), 200];
		let data = await vote_enqueuer_handler(newbody);
		return [data, 200];
	}
	return ['invalid request type', 200];
}

// Export the function
module.exports = { FunctionHandler };
