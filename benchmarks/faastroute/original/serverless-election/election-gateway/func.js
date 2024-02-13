const require('./rpc')
const axios = require("axios");

const function_handler = async (context) => {
	if (body['requestType'] ==  'get_results') {
		var data = '';
                await axios.post(process.env.ELECTION_GET_RESULTS, body)
                        .then( (response) => {
                                data = response.data;
                        });
		return data;
	}
	else if (body['requestType'] == 'vote') {
		var data = '';
		await axios.post(process.env.ELECTION_VOTE_ENQUEUER, body)
			.then( (response) => {
				data = response.data;
			});
		return data;
	}
	return 'invalid request type';
}

// Export the function
module.exports = { function_handler };
