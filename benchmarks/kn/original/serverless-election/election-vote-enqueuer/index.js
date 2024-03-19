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
const moment = require('moment');

const client = redis.createClient({url: process.env.REDIS_URL, password: process.env.REDIS_PASSWORD});

if ("LOGGING_NAME" in process.env) {
        const loggingClient = redis.createClient({url: process.env.LOGGING_URL, password: process.env.LOGGING_PASSWORD});
}

const handle = async (context, body) => {
	if ("LOGGING_NAME" in process.env) {
                await loggingClient.append(process.env.LOGGING_NAME, moment().format('MMMM Do YYYY, h:mm:ss a') + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "0" + "\n");
        }
	client.on('error', err => console.log('Redis Client Error', err));
        await client.connect();
	reply = await client.exists("voter-" + body['id']);
	if (reply == 1) {
		const g_val = await client.get("voter-" + body['id']);
		if (g_val != "Not Voted") {
			return {"isBase64Encoded": false, "statusCode": 409, "body": {"success": false, "message": (body['id'] + " already submitted a vote.")}};
		}
		else {
			data = '';
			newBody = body;
			if ("LOGGING_NAME" in process.env) {
        		        await loggingClient.append(process.env.LOGGING_NAME, moment().format('MMMM Do YYYY, h:mm:ss a') + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "1" + "\n");
        		}
			await axios.post(process.env.ELECTION_VOTE_PROCESSOR, newBody)
				.then( (response) => {
                                	data = response.data;
				});
			if ("LOGGING_NAME" in process.env) {
         		        await loggingClient.append(process.env.LOGGING_NAME, moment().format('MMMM Do YYYY, h:mm:ss a') + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "2" + "\n");
        		}
			return {"isBase64Encoded": false, "statusCode": 201, "body": {"success": true, "message": ("Vote " + body['id'] + " registered")}};
		}
	}
	if ("LOGGING_NAME" in process.env) {
                await loggingClient.append(process.env.LOGGING_NAME, moment().format('MMMM Do YYYY, h:mm:ss a') + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "3" + "\n");
        }
	return {"isBase64Encoded": false, "statusCode": 404, "body": {"success": false, "message": ("This voter id does not exist: " + body['id'])}};
}

// Export the function
module.exports = { handle };

