# Workflow Evaluation Definitions

This directory contains copies of used workflow evaluation definitions to be used as inputs to the Evaluator (see base for details). A evaluation definition consists of the workflow definitions (see workflow-definitions), any extra targets and their payloads, and the invocation count and each workflow concurrency to run the evaluation at. Listed below is an explanation of the format for the workflow evaluation definitions such that one can be constructed for testing and development.

1. Workflow evaluation definitions utilize the WorkflowStruct protobuf object for parsing (see base for explanation), thus all named variables must match (capitalization sensitive).
2. Evaluation definition files sent to HTTP must adhere to JSON formatting, otherwise the Evaluation Start calls will be interpreted as malformed.

Example:

	{
		"WorkflowConcurrency": [8, 16, 20, 40, 80, 100], // list of workflow concurrencies (i.e. parallel workflow invocations) to be testing every benchmark with
		"Invocations": 1000, // total number of invocations to be executed per benchmark, workflow concurrency and phase (warm start and cold start) pairing
		"ExtraTargets": { // list of extra target URLs to be evaluated using the same scheme as the MacroPod variants (i.e. Knative and other frameworks)
			"kn-election-original": "http://election-gateway.knative-functions.192.168.56.20.sslip.io",
			"kn-election-unified": "http://election-unified.knative-functions.192.168.56.20.sslip.io"
		},
		"ExtraTargetsPayload": { // extra target payloads; list of PayloadStructs used exactly like what is used in workflow definition payloads
			"kn-election-original": {
				"Type": "JSON",
				"JSON": {"requestType": "vote", "id": "1234", "state": "TX", "candidate": "John_Smith"}
			},
			"kn-election-unified": {
				"Type": "JSON",
				"JSON": {"requestType": "vote", "id": "1234", "state": "TX", "candidate": "John_Smith"}
			}
		},
		"Workflows": { // workflow definitions that are to be evaluated; note all names must be unique from each other and to the Ingress Controller
			"macropod-election-agg": {
				"Name": "election",
				"Functions": {
					"ELECTION_GATEWAY": {
						"Registry": "arearnest7/election-gateway:latest",
						"Endpoints": ["ELECTION_GET_RESULTS", "ELECTION_VOTE_ENQUEUER"]
					},
					"ELECTION_GET_RESULTS": {
						"Registry": "arearnest7/election-get-results:latest",
						"Envs": {"REDIS_URL": "127.0.0.1", "REDIS_PASSWORD": "password"}
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
				"Config": {
					"Deployment": "macropod",
					"Communication": "direct",
					"Aggregation": "agg"
				},
				"Payload": {
					"Type": "JSON",
					"JSON": {"requestType": "vote", "id": "1234", "state": "TX", "candidate": "John_Smith"}
				}
			}
		}
	}