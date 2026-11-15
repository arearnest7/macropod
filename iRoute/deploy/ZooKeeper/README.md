# Deploying Zookeeper on Kubernetes

This guide provides the steps to deploy a Zookeeper cluster on Kubernetes using configuration files.

## Prerequisites
- A running Kubernetes cluster
- `kubectl` installed and configured to access your cluster

## Deployment Steps

1. **Install the Local Path Provisioner**

   ```bash
   kubectl apply -f local-path-storage.yaml
   ```

2. **Apply the Zookeeper configuration**

   ```bash
   kubectl apply -f zookeeper-config.yaml
   ```

3. **Create the headless service**

   ```bash
   kubectl apply -f zookeeper-headless-svc.yaml
   ```

4. **Deploy the Zookeeper StatefulSet**

   ```bash
   kubectl apply -f zookeeper-statefulset.yaml
   ```

## Verify the Deployment
Check if the Zookeeper pods are running:
   ```bash
   kubectl get pods -l app=zookeeper
   ```
You should see the list of Zookeeper pods with their status as `Running`.