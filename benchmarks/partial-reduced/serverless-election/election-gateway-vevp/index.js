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

const state_list = ['AK', 'AL', 'AR', 'AZ', 'CA', 'CO', 'CT', 'DC', 'DE', 'FL', 'GA', 'HI', 'IA', 'ID'
, 'IL', 'IN', 'KS', 'KY', 'LA', 'MA', 'MD', 'ME', 'MI', 'MN', 'MO', 'MS', 'MT', 'NC', 'ND', 'NE', 'NH'
, 'NJ', 'NM', 'NV', 'NY', 'OH', 'OK', 'OR', 'PA', 'RI', 'SC', 'SD', 'TN', 'TX', 'U'];

const vote_processor_handler = async (body) => {
        client.set("voter-" + body['id'], body);

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

const vote_enqueuer_handler = async (body) => {
        client.exists("voter-" + body['id'] , function(err, reply) {
                if (reply === 1) {
                        const g_val = client.get("voter-" + body['id']);
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

const handle = async (context, body) => {
        var opt = {
                host: process.env.ELECTION_GET_RESULTS_PARTIAL,
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
                let data = await vote_enqueuer_handler(body);
                return data;
        }
        return 'invalid request type';
}

// Export the function
module.exports = { handle };
