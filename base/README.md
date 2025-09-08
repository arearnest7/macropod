# MacroPod Base Images

![ ](/home/arearnest7/macropod/base/macropod_diagram.png  "MacroPod Diagram")

1. Ingress Controller (macropod-ingress): The ingress controller is the entry point of MacroPod. It accepts and forwards the requests to the workflow's entry functions if the service is available. However, if the service is not available, it triggers the creation of new pods. For load balancing, the ingress controller forwards the request to the replica that has the least number of active requests.

	inputs (environment variables):

		SERVICE_PORT: string defining gRPC service binding; default 8000
		HTTP_PORT: string defining HTTP service binding; default 9000
		DEPLOYER_ADDRESS: string IP:port for the deployer gRPC service; default 127.0.0.1:8002
		NAMESPACE: default namespace if workflow definition value not set; default macropod-functions
		TTL: default time to live if workflow definition value not set; default 180 (180s)
		DEPLOYMENT: default deployment method if workflow definition value not set; default macropod
		COMMUNICATION: default communication method if workflow definition value not set; default direct
		AGGREGATION: default aggregation method if workflow definition value not set; default agg
		TARGET_CONCURRENCY: default target concurrency if workflow definition value not set; default -1 (unbounded)
		DEBUG: debug level (higher is more verbose); default 0 (off)

 2. Deployer Controller (macropod-deployer): The deployment controller runs in the background and deploys new pods when triggered by the ingress controller. The MacroPod deployments are deployed such that during burst requests, each worker node has at most one replica of the pod. This strategy minimizes the resources consumed by a single replica during creation, execution, and deletion in a worker node. 

	inputs (environment variables):
	
		SERVICE_PORT: string defining gRPC service binding; default 8000
		HTTP_PORT: string defining HTTP service binding; default 9000
		INGRESS_ADDRESS: string IP:port for the ingress gRPC service; default 127.0.0.1:8001
		LOGGER_ADDRESS: string IP:port for the logger gRPC service; default is nil
		NAMESPACE: default namespace if workflow definition value not set; default macropod-functions
		TTL: default time to live if workflow definition value not set; default 180 (180s)
		DEPLOYMENT: default deployment method if workflow definition value not set; default macropod
		COMMUNICATION: default communication method if workflow definition value not set; default direct
		AGGREGATION: default aggregation method if workflow definition value not set; default agg
		TARGET_CONCURRENCY: default target concurrency if workflow definition value not set; default -1 (unbounded)
		DEBUG: debug level (higher is more verbose); default 0 (off)

3. Function Logger (macropod-logger): This logger service is designed to unify and timestamp all prints, error and notable timestamps (invocations and RPCs) for a workflow by function. This information can be captured out to file for easy debugging of functions without Kubernetes logs being reset due to failure or restart. For the sake of collection the logger was disabled to remove any potential latency variance between MacroPod and Knative versions.

	inputs (environment variables):
	
		SERVICE_PORT: string defining gRPC service binding; default 8000
		HTTP_PORT: string defining HTTP service binding; default 9000
		TIMESTAMP_DIR: directory for storing all timestamp data; default /app/timestamp/
		ERROR_DIR: directory for storing all error data; default /app/error/
		PRINT_DIR: directory for storing all print data; default /app/print/

4. Metrics Server (macropod-metrics): This service is a Kubernetes service deployed onto all worker nodes to collect system-wide information such as CPU and memory utilization via system process outputs. Doing so the evaluator is able to get cluster wide system information during collections.

	inputs (environment variables):
	
		SERVICE_PORT: string defining gRPC service binding; default 8000
		HTTP_PORT: string defining HTTP service binding; default 9000

5. Evaluator (macropod-eval): The evaluator takes in an evaluation definition that defines which workflows, including their configurations and their evaluation payload, the workflow concurrencies and the number of invocations per workflow concurrency is to be collected. Furthermore the definition can be specified with extra targets as Knative to test as well. The evaluation process itself collects resource utilization from the Metrics Collectors (one on each worker) while it invokes the specified number of cold starts and warm starts, collecting the E2E latency at the same time. This is designed to allow for batch evaluation with specified configurations.

	inputs (environment variables):
	
		SERVICE_PORT: string defining gRPC service binding; default 8000
		HTTP_PORT: string defining HTTP service binding; default 9000
		INGRESS_ADDRESS: string IP:port for the ingress gRPC service; default 127.0.0.1:8000
		INGRESS_ADDRESS_HTTP: string IP:port for the ingress HTTP service; default 127.0.0.1:9000
		METRICS_DIR: directory for storing all metrics results; default /app/metrics/
		LATENCY_DIR: directory for storing all latency results; default /app/latency/
		SUMMARY_DIR: directory for storing all summary results; default /app/summary/
		WORKER_NODES: array of all worker node IPs as a string (for pulling stats); default "192.168.56.21 192.168.56.22 192.168.56.23 192.168.56.24"
		DEBUG: debug level (higher is more verbose); default 0 (off)
		TTL: time to live for all evaluation requests; default 360 (360s)

6. Go Function (macropod-go): This is the Golang version of the workflow function handler. This template handles any Golang function built with it to interact with MacroPod. It is also designed to work with HTTP calls as well in case of need for localized testing without other components.

	inputs (environment variables):
	
		SERVICE_PORT: string defining gRPC service binding; default 5000
		HTTP_PORT: string defining HTTP service binding; default 6000
		WORKFLOW: string name of workflow the function belongs to
		FUNCTION: string name of the function
		INGRESS: string representing the IP:port of the ingress such that gateway communication can be used if enabled
		LOGGER: string IP:port of the logger if logging is enabled; empty if logging is disabled
		COMM_TYPE: string name of the communication type (refer to the deployer documentation)

7. NodeJS Function (macropod-node): This is the NodeJS version of the workflow function handler. This template handles any NodeJS function built with it to interact with MacroPod. It is also designed to work with HTTP calls as well in case of need for localized testing without other components.

	inputs (environment variables):
	
		SERVICE_PORT: string defining gRPC service binding; default 5000
		HTTP_PORT: string defining HTTP service binding; default 6000
		WORKFLOW: string name of workflow the function belongs to
		FUNCTION: string name of the function
		INGRESS: string representing the IP:port of the ingress such that gateway communication can be used if enabled
		LOGGER: string IP:port of the logger if logging is enabled; empty if logging is disabled
		COMM_TYPE: string name of the communication type (refer to the deployer documentation)

8. Python Function (macropod-python): This is the Python3 version of the workflow function handler. This template handles any Python3 function built with it to interact with MacroPod. It is also designed to work with HTTP calls as well in case of need for localized testing without other components.

	inputs (environment variables):
	
		SERVICE_PORT: string defining gRPC service binding; default 5000
		HTTP_PORT: string defining HTTP service binding; default 6000
		WORKFLOW: string name of workflow the function belongs to
		FUNCTION: string name of the function
		INGRESS: string representing the IP:port of the ingress such that gateway communication can be used if enabled
		LOGGER: string IP:port of the logger if logging is enabled; empty if logging is disabled
		COMM_TYPE: string name of the communication type (refer to the deployer documentation)

9. Protobuf: This is the gRPC protocol used by all MacroPod functions. Listed below is the latest copy of the protobuf structure annotated.

		syntax = "proto3";
		import "google/protobuf/struct.proto";

		package macropod;
		option go_package = "./macropod_pb";

		message MacroPodRequest {													// general purpose request structure; used for all requests that do not need a specific request format; used for MacroPod function payloads
		    optional string Text = 1;												// text/string payload
		    optional google.protobuf.Struct JSON = 2;								// general purpose JSON structure; does not have stated format, but requires it to comply with JSON formatting for it to not be nil
		    optional bytes Data = 3;												// byte array payload
		    optional string Workflow = 4;											// workflow name
		    optional string Function = 5;											// function name
		    optional string WorkflowID = 6;											// workflow invocation ID; used to track unique request chains for logging
		    optional int32 Depth = 7;												// workflow invocation depth at current request; used for logging
		    optional int32 Width = 8;												// workflow invocation count in parallel/fan-out call; used for logging
		    optional string Target = 9;												// string representing ingress path taken to get to location; used for logging and for general purpose string assignment for other requests
		}

		message MacroPodReply {														// general purpose reply structure; used for all requests that do not need a specific reply format
		    optional string Reply = 1;												// response string; usually a return message from call path signifying execution results
		    optional int32 Code = 2;												// response code; used exclusively to track success or failure of workflow invocation calls
		}

		message FunctionStruct {													// function definition structure; used as an intermediate value for defining a workflow
		    string Registry = 1;													// image registry and tag; this defines where a container image must be pulled from if not pre-existing in cluster
		    repeated string Endpoints = 2;											// list of downstream endpoints this function can use; they must match the naming scheme found in the WorkflowStruct definition
		    map<string, string> Envs = 3;											// any environment variables that the function needs to operate; please note that if the name overlaps with any of the names set in MacroPod function inputs (see above) they will be overwritten
		    map<string, string> Secrets = 4;										// any Kubernetes secrets that the function needs to operate
		}

		message ConfigStruct {														// configuration structure for workflows and the default handler on ingress/deployer controller; in terms of default, if workflow definition does not have the value set it will defer to the default for value
		    optional string Namespace = 1;											// namespace for all replica deployments
		    optional int32 TTL = 2;													// time-to-live for each function in seconds
		    optional string Deployment = 3;											// deployment method; signifies how the workflow functions will be deployed in cluster; for now only handles ["macropod"]
		    optional string Communication = 4;										// communication strategy; signifies routing patterns and expected destination paths for downstream functions; for now handles ["direct", "gateway"]
		    optional string Aggregation = 5;										// aggregation/disaggregation strategy; signifies how the functions of a workflow will be grouped and works in tandem with deployment to structure and layout replicas; currentl handles ["agg", "disagg", "dynamic"]
		    optional int32 TargetConcurrency = 6;									// signifies the maximum number of requests any given replica can be given by ingress controller; -1 means unbounded/infinite
		    optional int32 Debug = 7;												// debug level for ingress/deployer; is not used by workflow definition
		}

		message PayloadStruct {														// payload definition for evaluations
		    string Type = 1;														// sets which payload structure to use for evaluation; string must match the name of the structure
		    optional string Text = 2;												// text/string payload
		    optional google.protobuf.Struct JSON = 3;								// general purpose JSON structure payload
		    optional bytes Data = 4;												// bytes array payload
		}

		message WorkflowStruct {													// workflow definition structure; used to set/create a workflow in ingress controller and is a component of the EvalStruct for defining which benchmarks/workflows to run and how
		    string Name = 1;														// workflow name
		    map<string, FunctionStruct> Functions = 2;								// map of all workflow functions mapped by name
		    optional ConfigStruct Config = 3;										// configuration defining how to deploy, communicate and aggregate the workflow
		    optional PayloadStruct Payload = 4;										// payload to be used during evaluations of this workflow
		}

		message EvalStruct {														// evaluation definition; used to evaluate multiple benchmarks/workflows under MacroPod's Evaluator
		    repeated int32 WorkflowConcurrency = 1;									// list of all workflow concurrencies (i.e. number of workflow invocations concurrently running) to be tested with
		    int32 Invocations = 2;													// number of invocations to be used per benchmark/workflow, concurrency and phase combination
		    map<string, string> ExtraTargets = 3;									// map of all extra invocation targets; these will be handled the same as the MacroPod workflows, however will be invoked directly instead of contacting Ingress Controller
		    map<string, PayloadStruct> ExtraTargetsPayload = 4;						// extra invocation target payloads
		    map<string, WorkflowStruct> Workflows = 5;								// map of all workflow definitions to be used; please note provide a payload for each benchmark/workflow otherwise it will not work; name used in map is what will be set in output lines for that iteration of the benchmark if multiple exist
		}

		message MetricsStruct {														// system metrics from a metrics server; utilizes github.com/shirou/gopsutil to provide system stats; please consult repository or psutil man page
		    optional double Uptime = 1;
		    optional double LoadAvg1 = 2;
		    optional double LoadAvg5 = 3;
		    optional double LoadAvg15 = 4;
		    optional double CPUUsed = 5;
		    optional double CPUCountLogical = 6;
		    optional double CPUCountPhysical = 7;
		    optional double MemoryUsed = 8;
		    optional double MemoryAvailable = 9;
		    optional double MemoryTotal = 10;
		    optional double MemoryBuffers = 11;
		    optional double MemoryCached = 12;
		    optional double MemoryWriteBack = 13;
		    optional double MemoryDirty = 14;
		    optional double MemoryWriteBackTmp = 15;
		    optional double MemoryShared = 16;
		    optional double MemorySlab = 17;
		    optional double MemorySreclaimable = 18;
		    optional double MemorySunreclaim = 19;
		    optional double MemoryPageTables = 20;
		    optional double MemorySwapCached = 21;
		    optional double MemoryCommitLimit = 22;
		    optional double MemoryCommittedAS = 23;
		    optional double MemoryHighTotal = 24;
		    optional double MemoryHighFree = 25;
		    optional double MemoryLowTotal = 26;
		    optional double MemoryLowFree = 27;
		    optional double MemorySwapTotal = 28;
		    optional double MemorySwapFree = 29;
		    optional double MemoryMapped = 30;
		    optional double MemoryVmallocTotal = 31;
		    optional double MemoryVmallocUsed = 32;
		    optional double MemoryVmallocChunk = 33;
		    optional double MemoryHugePagesTotal = 34;
		    optional double MemoryHugePagesFree = 35;
		    optional double MemoryHugePagesRsvd = 36;
		    optional double MemoryHugePagesSurp = 37;
		    optional double MemoryHugePageSize = 38;
		    optional double MemoryAnonHugePages = 39;
		    optional double DiskUsed = 40;
		    optional double DiskFree = 41;
		    optional double DiskTotal = 42;
		    optional double DiskInodesUsed = 43;
		    optional double DiskInodesFree = 44;
		    optional double DiskInodesTotal = 45;
		    optional double NetworkBytesSent = 46;
		    optional double NetworkBytesRecv = 47;
		    optional double NetworkPacketsSent = 48;
		    optional double NetworkPacketsRecv = 49;
		    optional double NetworkErrin = 50;
		    optional double NetworkErrout = 51;
		    optional double NetworkDropin = 52;
		    optional double NetworkDropout = 53;
		    optional double NetworkFifoin = 54;
		    optional double NetworkFifoout = 55;
		}

		service MacroPodFunction {													// general purpose function handler
		    rpc Invoke(MacroPodRequest) returns (MacroPodReply);					// serve the function invoke and return the results
		}

		service MacroPodEval {														// Evaluator
		    rpc Eval(EvalStruct) returns (MacroPodReply);							// begins evaluation based on the evlaution definition provided, cycling over all benchmark, concurrency and phases variants
		    rpc EvalMetrics(MacroPodRequest) returns (MacroPodReply);				// returns a string of all contents from the metrics results for specified ID
		    rpc EvalLatency(MacroPodRequest) returns (MacroPodReply);				// returns a string of all contents from the latency results for specified ID
		    rpc EvalSummary(MacroPodRequest) returns (MacroPodReply);    			// returns a string of all contents from the summary results for specified ID
		}

		service MacroPodIngress {													// Ingress Controller
		    rpc Config(ConfigStruct) returns (MacroPodReply);						// sets global default configurations for all workflows
		    rpc WorkflowInvoke(MacroPodRequest) returns (MacroPodReply);			// invokes the requested workflow with payload; if a replica is not free or exists a deployment is created; if dynamic aggregation strategies used it will also trigger an update
		    rpc FunctionInvoke(MacroPodRequest) returns (MacroPodReply);			// for gateway communication patterns this is used to relay requests from upstream to downstream functions, returning results
		    rpc CreateWorkflow(WorkflowStruct) returns (MacroPodReply);				// create workflow entry based on workflow definition (WorkflowStruct)
		    rpc UpdateWorkflow(WorkflowStruct) returns (MacroPodReply);				// update workflow entry based on workflow definition (WorkflowStruct)
		    rpc DeleteWorkflow(MacroPodRequest) returns (MacroPodReply);			// delete workflow from cluster
		}

		service MacroPodDeployer {													// Deployer Controller
		    rpc Config(ConfigStruct) returns (MacroPodReply);						// sets global default configurations for all workflows
		    rpc CreateWorkflow(WorkflowStruct) returns (MacroPodReply);				// create workflow entry based on workflow definition (WorkflowStruct)
		    rpc UpdateWorkflow(WorkflowStruct) returns (MacroPodReply);				// update workflow entry based on workflow definition (WorkflowStruct)
		    rpc DeleteWorkflow(MacroPodRequest) returns (MacroPodReply);			// delete workflow from cluster
		    rpc UpdateDeployments(MacroPodRequest) returns (MacroPodReply);			// updates workflow definitions if any dynamic aggregation methods is enabled
		    rpc CreateDeployment(MacroPodRequest) returns (MacroPodReply);			// create a new workflow deployment based on current cluster conditions
		    rpc TTLDelete(MacroPodRequest) returns (MacroPodReply);					// deletes replica target based on ingress signaling TTL is up
		}

		service MacroPodLogger {													// Logger
		    rpc Timestamp(MacroPodRequest) returns (MacroPodReply);					// append to timestamp file for the workflow:function combination the sent message and invocation data
		    rpc Error(MacroPodRequest) returns (MacroPodReply);						// append to error file for the workflow:function combination the sent message and invocation data
		    rpc Print(MacroPodRequest) returns (MacroPodReply);						// append to print file for the workflow:function combination the sent message and invocation data
		    rpc GetTimestamp(MacroPodRequest) returns (MacroPodReply);				// retrieve all timestamp logs from the function and workflow specified
		    rpc GetError(MacroPodRequest) returns (MacroPodReply);					// retrieve all error logs from the function and workflow specified
		    rpc GetPrint(MacroPodRequest) returns (MacroPodReply);					// retrieve all print logs from the function and workflow specified
		}

		service MacroPodMetrics {													// Metrics Server
		    rpc GetMetrics(MacroPodRequest) returns (MetricsStruct);				// retrieve current reported metrics on server
		}