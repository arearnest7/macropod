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
const http = require("http");

const handle = async (context, body) => {
	var opt = {
		host: process.env.ELECTION_GET_RESULTS,
		port: 8080,
		method: 'GET',
		headers: {
			'Content-Type': 'application/json',
    			'Content-Length': body
		}
	};
	if (body['requestType'] ==  'get_results') {
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
		return data;
	}
	else if (body['requestType'] == 'vote') {
		opt['host'] = process.env.ELECTION_VOTE_ENQUEUER;
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
                return data;
	}
	return 'invalid request type';
}

// Export the function
module.exports = { handle };
