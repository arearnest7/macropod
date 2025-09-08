# Knative Deployment Scripts

This directory is dedicated to any Knative deployment scripts necessary for evaluation. This consists of both deployment and removal scripts.

1. All deployment scripts are given in the form deploy-{workflow name}-{aggregation model}.sh and are invoked as "./deploy-{workflow name}-{aggregation model}.sh {host IP}"
2. Make sure to deploy a Knative cluster like that of in the deploy.sh script in the main directory before usage
3. To remove a given workflow invoke the removal script "./remove-{workflow name}-{aggregation model}.sh"