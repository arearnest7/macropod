# Deploying SocialNetwork Database on Kubernetes

This guide provides instructions to deploy the SocialNetwork database on a Kubernetes cluster.

## Prerequisites

- A running Kubernetes cluster
- `kubectl` installed and configured

---

## 1. Create the Namespace

Apply the namespace manifest to create the `socialnetwork-db` namespace:

```bash
kubectl apply -f namespaces.yml
```

This will create the `socialnetwork-db` namespace for all database resources.

## 2. Deploy the Database Components

Apply the database manifests:

```bash
kubectl apply -f ./yaml/
```

This will deploy all necessary database resources, such as:
- `Redis`
- `MongoDB`
- `memcached`

## 3. Verify the Deployment

Check whether all database pods are running in the `socialnetwork-db` namespace:

```bash
kubectl get pods -n socialnetwork-db
```
All pods should show the status `Running` or `Completed`.