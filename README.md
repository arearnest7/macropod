# MacroPod: A Universal Serverless Workflow Aggregation and Unification Testbed

This framework is an open-source testbed/benchmarking suite for the purpose of fairly evaluating serverless computing workflow scheduling methodology in terms of deployment, communication and aggregation strategies. This framework provides:

1. A minimalist framework designed for testing/development and redesign for new methods of aggregation, communication and deployment as a set of container services
2. A set of container services for testing, including a general purpose evaluator, logger and metrics server
3. A series of common benchmarks designed in both MacroPod and Knative formats along with deployment scripts for both for testing and evaluation for said strategies

This work has been submitted to **TBA** as **TBA** and the paper is accessible at **TBA**.

Listed below are instructions for provided scripts. All further instructions/information are in their respective base directory.

scripts:

	build.sh
		purpose: build all images in benchmarks/macropod and store contents in registry using ID for user and the appropriate tag
		usage: ./build.sh {ID} {TAG}
		inputs:
			ID: user account for DockerHub (or account linked using docker)
			TAG: container image tag; default is latest
	deploy.sh
		purpose: deploy a K3s (lightweight K8s) cluster onto the master node and worker nodes with both MacroPod and Knative components; please note that Knative is exceedingly chunky under basic deployment and requires a master node with at least 4 CPU/vCPU and 8GB RAM
		usage: ./deploy.sh {username} {password} {interface} {worker_nodes}
		inputs:
			username: user account across all nodes that has sudo privileges
			password: user account password
			interface: network interface on all nodes that you want the K3s networking to path through; if you do not know choose one that is not "lo"
			worker_nodes: a string list of all worker node IPs; please make sure your master node can reach all worker nodes
	eval.sh
		purpose: begins an evaluation using the evaluation definitions set in the file, by default election, hotel, sentiment and video from eval-definitions, and then writes the contents of the collection under the provided directory
		usage: ./eval.sh {HOST} {DIR}
		inputs:
			HOST: host IP for the Knative cluster; if you are not running Knative collections please comment out the deploy and remove scripts within the eval.sh script
			DIR: collection results directory; will write the metrics, latency and summary results from evaluation per evaluation definition provided
	remove.sh
		purpose: remove the K3s cluster from host and worker nodes
		usage: ./remove.sh {username} {password} {worker_nodes}
		inputs:
			username: user account across all nodes that has sudo privileges
			password: user account password
			worker_nodes: a string list of all worker node IPs; please make sure your master node can reach all worker nodes