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

fs.readFile('/etc/secret-volume/election-vote-processor', 'utf8', function(err, election-vote-processor) {
    if (err) throw err;
    console.log(election-vote-processor)
});
fs.readFile('/etc/secret-volume/redis-url', 'utf8', function(err, redis-url) {
    if (err) throw err;
    console.log(redis-url)
});

const client = redis.createClient({url: redis-url});

const handle = async (context, body) => {
	client.exists("voter-" + body['id'] , function(err, reply) {
		if (reply === 1) {
			const g_val = client.get(state_results[i]);
			if (g_val !== null) {
				return {"isBase64Encoded": false, "statusCode": 409, "body": {"success": false, "message": (body['id'] + "already submitted a vote.")}};
			}
			else {
				var opt = {
                			host: process.env.VOTE_PROC,
                			port: 80,
                			method: 'GET',
                			headers: {
                        			'Content-Type': 'application/json',
                        			'Content-Length': body
                			}
        			};
				let data = '';
                		http.get(opt, (res) => {
                        		res.on('data', (chunk) => {
                                		data += chunk;
                        		});
                        		res.on('end', () => {
                                		console.log('Body', JSON.parse(data))
                        		});
                		}).on("error", (err) => {
                        		console.log("Error: ", err)
                		}).end();
				return {"isBase64Encoded": false, "statusCode": 201, "body": {"success": true, "message": ("Vote " + body['id'] + " registered")}};
			}
		}
		return {"isBase64Encoded": false, "statusCode": 404, "body": {"success": false, "message": ("This voter id does not exist: " + body['id'])}};
	});
}

// Export the function
module.exports = { handle };

