const rpc = require('./rpc')
const redis = require('redis');

//const client = redis.createClient({url: process.env.REDIS_URL, password: process.env.REDIS_PASSWORD});

async function FunctionHandler(context) {
	//client.on('error', err => console.log('Redis Client Error', err));
        //await client.connect();
	//await client.set("voter-" + body['id'], JSON.stringify(body));
	var body = context.JSON;
	var state = body['state'];
	var candidate = body['candidate'];

	reply = 1; //await client.exists("election-results-" + state + "-" + candidate);
        if (reply == 1) {
                var cnt = 1; //parseInt(await client.get("election-results-" + state + "-" + candidate));
                cnt = cnt + 1;
                //await client.set("election-results-" + state + "-" + candidate, cnt.toString());
        }
        //else {
        //        await client.set("election-results-" + state + "-" + candidate, "1");
        //}
        return ["success", 200];
}

// Export the function
module.exports = { FunctionHandler };
