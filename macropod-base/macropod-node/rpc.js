const exec = require('child_process').execSync;
const moment = require('moment')
const axios = require("axios");

async function RPC(context, dest, payloads) {
	await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + context.WorkflowId + "," + context.Depth.toString() + "," + context.Width.toString() + "," + context.RequestType + "," + "rpc_start");
	var tl = new Array();
	var pv_paths = new Array();
	var request_type = "gg";
	var i = 0;
	for (let payload of payloads) {
		var request = {
			Data: payload,
			WorkflowId: context.WorkflowId,
			Depth: context.Depth+1,
			Width: i,
			RequestType: request_type
		};
		await axios.post("http://" + dest, request)
			.then( (response) => {
				var data = response.data;
				console.log(response);
				tl.push(data);
			});
		i += 1;
	}
	var results = new Array();
	for (let t of tl) {
		var reply = t;
		results.push(reply);
	}
        await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + context.WorkflowId + "," + context.Depth.toString() + "," + context.Width.toString() + "," + context.RequestType + "," + "rpc_end");
	return results;
}

// Export the function
module.exports = { RPC };
