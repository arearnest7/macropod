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
const axios = require("axios");
const moment = require('moment');
const redis = require('redis');

const handle = async (context, body) => {
        console.log(moment().format('MMMM Do YYYY, h:mm:sss a') + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "0" + "\n");
	if (body['requestType'] ==  'get_results') {
		var data = '';
        	console.log(moment().format('MMMM Do YYYY, h:mm:sss a') + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "1" + "\n");
                await axios.post(process.env.ELECTION_GET_RESULTS, body)
                        .then( (response) => {
                                data = response.data;
                        });
        	console.log(moment().format('MMMM Do YYYY, h:mm:sss a') + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "2" + "\n");
		return data;
	}
	else if (body['requestType'] == 'vote') {
		var data = '';
        	console.log(moment().format('MMMM Do YYYY, h:mm:sss a') + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "3" + "\n");
		await axios.post(process.env.ELECTION_VOTE_ENQUEUER, body)
			.then( (response) => {
				data = response.data;
			});
        	console.log(moment().format('MMMM Do YYYY, h:mm:sss a') + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "4" + "\n");
		return data;
	}
        console.log(moment().format('MMMM Do YYYY, h:mm:sss a') + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "5" + "\n");
	return 'invalid request type';
}

// Export the function
module.exports = { handle };
