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

//const client = redis.createClient({url: process.env.REDIS_URL, password: process.env.REDIS_PASSWORD});

const handle = async (context, body) => {
        var newbody = body;
	//client.on('error', err => console.log('Redis Client Error', err));
        //await client.connect();
	reply = 1; //await client.exists("voter-" + body['id']);
	if (reply == 1) {
		const g_val = "Not Voted"; //await client.get("voter-" + body['id']);
		if (g_val != "Not Voted") {
			return {"isBase64Encoded": false, "statusCode": 409, "body": {"success": false, "message": (body['id'] + " already submitted a vote.")}};
		}
		else {
			data = '';
			await axios.post(process.env.ELECTION_VOTE_PROCESSOR, newbody)
				.then( (response) => {
                                	data = response.data;
				});
			return {"isBase64Encoded": false, "statusCode": 201, "body": {"success": true, "message": ("Vote " + body['id'] + " registered")}};
		}
	}
	return {"isBase64Encoded": false, "statusCode": 404, "body": {"success": false, "message": ("This voter id does not exist: " + body['id'])}};
}

// Export the function
module.exports = { handle };

