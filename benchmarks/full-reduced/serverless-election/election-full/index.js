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

const fs = require('fs')
const redis = require('redis');
const http = require('http');

fs.readFile('/etc/secret-volume/redis-url', 'utf8', function(err, redis-url) {
    if (err) throw err;
    console.log(redis-url)
});

const client = redis.createClient({url: redis-url});

const state_list = ['AK', 'AL', 'AR', 'AZ', 'CA', 'CO', 'CT', 'DC', 'DE', 'FL', 'GA', 'HI', 'IA', 'ID'
, 'IL', 'IN', 'KS', 'KY', 'LA', 'MA', 'MD', 'ME', 'MI', 'MN', 'MO', 'MS', 'MT', 'NC', 'ND', 'NE', 'NH'
, 'NJ', 'NM', 'NV', 'NY', 'OH', 'OK', 'OR', 'PA', 'RI', 'SC', 'SD', 'TN', 'TX', 'U'];

const vote_processor_handler = async (body) => {
        client.set(body['id'], body);

        var state = body['state'];
        var candidate = body['candidate'];

        client.exists("election-results-" + state + "-" + candidate, function(err, reply) {
                if (reply === 1) {
                        var cnt = parseInt(client.get("election-results-" + state + "-" + candidate));                        cnt = cnt + 1;
                        client.set("election-results-" + state + "-" + candidate, cnt.toString());
                }
                else {
                        client.set("election-results-" + state + "-" + candidate, "1");
                }

        });
        return "success";
}

const vote_enqueuer_handler = async (body) => {
        client.exists("voter-" + body['id'] , function(err, reply) {
                if (reply === 1) {
                        const g_val = client.get(state_results[i]);
                        if (g_val !== null) {
                                return {"isBase64Encoded": false, "statusCode": 409, "body": {"success": false, "message": (body['id'] + "already submitted a vote.")}};
                        }
                        else {
                                let data = vote_processor_handler(body);
                                return {"isBase64Encoded": false, "statusCode": 201, "body": {"success": true, "message": ("Vote " + body['id'] + " registered")}};
                        }
                }
                return {"isBase64Encoded": false, "statusCode": 404, "body": {"success": false, "message": ("This voter id does not exist: " + body['id'])}};
        });
}

const get_results_handler = async (body) => {
        var results = [];
        for (var state in state_list) {
                client.keys('election-results-' + state + '-*', function (err, state_results) {
                        if (err) {
                                return console.log(err);
                        }
                        for(var i = 0, len = keys.length; i < len; i++) {
                                console.log(keys[i]);
                        }
                });
                var total_count = {"total": 0};
                for (var i = 0; i < state_results; i++) {
                        const cnt = await client.get(state_results[i]);
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
        if (body['requestType'] ==  'get_results') {
                let data = await get_results_handler(body);
                return data;
        }
        else if (body['requestType'] == 'vote') {
                let data = await vote_enqueuer_handler(body);
		return data;
        }
        return 'invalid request type';
}

// Export the function
module.exports = { handle };
