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
const redis = require('redis');
const exec = require('child_process').execSync;
const moment = require('moment');

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
        console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + workflow_depth.toString() + "," + workflow_width.toString() + "," + "HTTP" + "," + "0");
	if (body['requestType'] ==  'get_results') {
		var data = '';
        	console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + workflow_depth.toString() + "," + workflow_width.toString() + "," + "HTTP" + "," + "1");
                await axios.post(process.env.ELECTION_GET_RESULTS, newbody)
                        .then( (response) => {
                                data = response.data;
                        });
                console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + workflow_depth.toString() + "," + workflow_width.toString() + "," + "HTTP" + "," + "2");
		return data;
	}
	else if (body['requestType'] == 'vote') {
		var data = '';
                console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + workflow_depth.toString() + "," + workflow_width.toString() + "," + "HTTP" + "," + "3");
		await axios.post(process.env.ELECTION_VOTE_ENQUEUER, newbody)
			.then( (response) => {
				data = response.data;
			});
                console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + workflow_depth.toString() + "," + workflow_width.toString() + "," + "HTTP" + "," + "4");
		return data;
	}
        console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + workflow_depth.toString() + "," + workflow_width.toString() + "," + "HTTP" + "," + "5");
	return 'invalid request type';
}

// Export the function
module.exports = { handle };
