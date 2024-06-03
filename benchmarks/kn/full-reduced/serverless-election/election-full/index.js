/**
 * Your HTTP handling function, invoked with each request. This is an example
 * function that echoes its input to the caller, and returns an error if
 * the incoming request is something other than an HTTP POST or GET.
 *
 * In can be invoked with 'func invoke'
 * It can be tested with 'npm test'
 *
 * @param {Context} context a context object.
 * @param {object} context.body the request body if any
 * @param {object} context.query the query string deserialized as an object, if any
 * @param {object} context.log logging object with methods for 'info', 'warn', 'error', etc.
 * @param {object} context.headers the HTTP request headers
 * @param {string} context.method the HTTP request method
 * @param {string} context.httpVersion the HTTP protocol version
 * See: https://github.com/knative/func/blob/main/docs/function-developers/nodejs.md#the-context-object
 */

const redis = require('redis');
const axios = require('axios');
const exec = require('child_process').execSync;
const moment = require('moment');

const client = redis.createClient({url: 'redis://' + process.env.REDIS_URL, password: process.env.REDIS_PASSWORD});

const state_list = ['AK', 'AL', 'AR', 'AZ', 'CA', 'CO', 'CT', 'DC', 'DE', 'FL', 'GA', 'HI', 'IA', 'ID'
, 'IL', 'IN', 'KS', 'KY', 'LA', 'MA', 'MD', 'ME', 'MI', 'MN', 'MO', 'MS', 'MT', 'NC', 'ND', 'NE', 'NH'
, 'NJ', 'NM', 'NV', 'NY', 'OH', 'OK', 'OR', 'PA', 'RI', 'SC', 'SD', 'TN', 'TX', 'U'];

const vote_processor_handler = async (body) => {
        await client.set("voter-" + body['id'], JSON.stringify(body));

        var state = body['state'];
        var candidate = body['candidate'];

	reply = await client.exists("election-results-" + state + "-" + candidate);
        if (reply == 1) {
        	var cnt = parseInt(await client.get("election-results-" + state + "-" + candidate));
		cnt = cnt + 1;
		await client.set("election-results-" + state + "-" + candidate, cnt.toString());
	}
	else {
		await client.set("election-results-" + state + "-" + candidate, "1");
	}
        return "success";
}

const vote_enqueuer_handler = async (body) => {
	reply = await client.exists("voter-" + body['id']);
	if (reply == 1) {
        	const g_val = await client.get("voter-" + body['id']);
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
                state_results = await client.keys('election-results-' + state + '-*');
                var total_count = {"total": 0};
                for (var i = 0; i < state_results.length; i++) {
                        const cnt = await parseInt(client.get(state_results[i]));
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

const handle = async (context, body) => {
	var workflow_id = Math.floor(Math.random() * 10000000).toString();
        var workflow_depth = 0;
        var workflow_width = 0;
        var newbody = body;
        if ("workflow_id" in body) {
                workflow_id = body["workflow_id"];
                workflow_depth = body["workflow_depth"] + 1;
                workflow_width = body["workflow_width"];
        } else {
                newbody["workflow_id"] = workflow_id;
                newbody["workflow_depth"] = workflow_depth;
                newbody["workflow_width"] = workflow_width;
        }
	console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + workflow_depth.toString() + "," + workflow_width.toString() + "," + "HTTP" + "," + "0" + "\n");
	client.on('error', err => console.log('Redis Client Error', err));
	await client.connect();
        if (body['requestType'] ==  'get_results') {
                let data = await get_results_handler(newbody);
        	console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + workflow_depth.toString() + "," + workflow_width.toString() + "," + "HTTP" + "," + "1" + "\n");
                return data;
        }
        else if (body['requestType'] == 'vote') {
                let data = await vote_enqueuer_handler(newbody);
                console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + workflow_depth.toString() + "," + workflow_width.toString() + "," + "HTTP" + "," + "2" + "\n");
		return data;
        }
        console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + workflow_depth.toString() + "," + workflow_width.toString() + "," + "HTTP" + "," + "3" + "\n");
        return 'invalid request type';
}

// Export the function
module.exports = { handle };

