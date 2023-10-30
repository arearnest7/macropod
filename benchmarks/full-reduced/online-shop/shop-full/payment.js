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
const path = require('path');
const grpc = require('@grpc/grpc-js');
const pino = require('pino');
const protoLoader = require('@grpc/proto-loader');

const charge = require('./charge');

const logger = pino({
    name: 'paymentservice-server',
    messageKey: 'message',
    formatters: {
        level (logLevelString, logLevelNum) {
            return { severity: logLevelString }
        }
    }
});

static chargeServiceHandler(request) {
	try {
            logger.info(`PaymentService#Charge invoked with request ${JSON.stringify(request)}`);
            const response = charge(request);
            callback(null, response);
        } catch (err) {
            console.warn(err);
            callback(err);
        }
}

function main(body) {
    return chargeServiceHandler(body);
}

if (require.main == module) {
    main(process.argv[2])
}
