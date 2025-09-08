# Workflow Definitions

This directory contains copies of used workflow definitions to be used as inputs for Ingress Controller (see base for details). A workflow definition contains the name of the workflow, its functions and their configuration (optional). Listed below is an explanation of the format for the workflow definitions such that one can be constructed for testing and development.

1. Workflow definitions utilize the WorkflowStruct protobuf object for parsing (see base for explanation), thus all named variables must match (capitalization sensitive).
2. Workflow definition files sent to HTTP must adhere to JSON formatting, otherwise the workflow create and update calls will be interpreted as malformed.

Example:

    {
		"Name": "election", // Workflow Name; unique name identifier for the Ingress Controller, so choose name wisely for multi-tenancy
		"Functions": { // list of functions that define the workflow
			"ELECTION_GATEWAY": { 
			// function, note that all function names are in all caps and spaces are given underscores. 
			// The names will be converted to this format internally. 
			// Endpoint names must match the name of the function in the definition.
			    "Registry": "arearnest7/election-gateway:latest", // registry name and tag for the function image; must be accessible/knowable by the K8s cluster
			    "Endpoints": ["ELECTION_GET_RESULTS", "ELECTION_VOTE_ENQUEUER"] // list of endpoints used by this function. If provided the function can call the listed functions as downstreams.
			},
			"ELECTION_GET_RESULTS": {
				"Registry": "arearnest7/election-get-results:latest",
				"Envs": {"REDIS_URL": "127.0.0.1", "REDIS_PASSWORD": "password"} // any environmental variables that need to be passed onto the function
			},
			"ELECTION_VOTE_ENQUEUER": {
				"Registry": "arearnest7/election-vote-enqueuer:latest",
				"Endpoints": ["ELECTION_VOTE_PROCESSOR"],
				"Envs": {"REDIS_URL": "127.0.0.1", "REDIS_PASSWORD": "password"}
			},
			"ELECTION_VOTE_PROCESSOR": {
			    "Registry": "arearnest7/election-vote-processor:latest",
			    "Envs": {"REDIS_URL": "127.0.0.1", "REDIS_PASSWORD": "password"}
			}
		},
		"Config": { // ConfigStruct that defines deployment, communication and aggregation strategies; see base for full list (optional)
			"Deployment": "macropod",
			"Communication": "direct",
			"Aggregation": "agg"
		},
		"Payload": { // payload to be used for evaluations (optional)
			"Type": "JSON", // payload type; can be Text/JSON/Data
			"JSON": {"requestType": "vote", "id": "1234", "state": "TX", "candidate": "John_Smith"} // JSON example
		}
    }