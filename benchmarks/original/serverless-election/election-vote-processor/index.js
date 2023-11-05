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

const client = redis.createClient({url: process.env.REDIS_URL});

const handle = async (context, body) => {
	client.set(body['id'], body);

	var state = body['state'];
	var candidate = body['candidate'];

	client.exists("election-results-" + state + "-" + candidate, function(err, reply) {
                if (reply === 1) {
                	var cnt = parseInt(client.get("election-results-" + state + "-" + candidate));
			cnt = cnt + 1;
			client.set("election-results-" + state + "-" + candidate, cnt.toString());
		}
		else {
			client.set("election-results-" + state + "-" + candidate, "1");
		}

	});
	return "success";
}

// Export the function
module.exports = { handle };
