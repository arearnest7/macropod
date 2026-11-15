# Deploying k8s-rdma-shared-dev-plugin on Kubernetes

This guide explains how to deploy the `k8s-rdma-shared-dev-plugin` to enable RDMA shared device support in Kubernetes.

---

## 1. Prerequisites

Before deployment, ensure the following:

- Kubernetes cluster is running
- `kubectl` configured and connected to your cluster
- RDMA-capable NICs installed on nodes
- Necessary RDMA libraries and drivers, e.g., `rdma-core`, installed on all nodes

---

## 2. Identify RDMA Devices and Corresponding Network Interfaces

Verify which network interfaces are associated with your RDMA devices.  

You can use one of the following methods on each node:

```bash
ibdev2netdev
```

Example output:
```
mlx5_0 port 1 ==> ens5f0np0 (Down)
mlx5_1 port 1 ==> ens5f1np1 (Up)
```

Record the mapping between RDMA devices (e.g., `mlx5_1`) and their network interfaces (e.g., `ens5f1np1`). This information is required for correct plugin configuration.

---


## 3. Deploy the RDMA Plugin

Replace `ens5f1np1` in `configmap.yaml` with the actual network interface(s) you identified in Step 2.

Deploy the plugin using the provided DaemonSet YAML:

```bash
kubectl apply -k .
```

This will deploy the plugin on all eligible nodes in the cluster.

---

## 4. Verify Deployment

Check that the DaemonSet pods are running:

```bash
kubectl get pods -n kube-system | grep rdma
```

All pods should be in `Running` state.

Also, check the logs for any errors:

```bash
kubectl logs -n kube-system <rdma-plugin-pod-name>
```

## 5. Verify RDMA Device Availability

After the plugin is deployed, RDMA devices should appear on the nodes:

```bash
kubectl get nodes -o json | jq '.items[].status.allocatable'
```

You should see RDMA-related resources, for example: `rdma/hca_shared_devices_a:0`.

---

## 6. Optional: Cleanup

To remove the RDMA plugin from the cluster:

```bash
kubectl delete -f configmap.yaml
```