# Deploying OpenFaaS on Kubernetes

This guide provides instructions to deploy OpenFaaS on Kubernetes using YAML manifests and `kubectl`.

## Prerequisites

- A running Kubernetes cluster
- `kubectl` installed and configured

---

## 1. Install the OpenFaaS CLI

Download and install the OpenFaaS CLI:

```bash
curl -sSL https://cli.openfaas.com | sudo -E sh
```

Verify the installation:
```bash
faas-cli version
```

## 2. Create Namespaces
Apply the namespace manifests:
```bash
kubectl apply -f ./namespaces.yml
```
This will create the following namespaces:
- `openfaas` (system components)

- `openfaas-fn` (user functions)

## 3. Create Basic Authentication Secret

Set the password for the admin user:

```bash
export PASSWORD=<your_password_here>
```

Create the Kubernetes secret in the `openfaas` namespace:
```bash
kubectl -n openfaas create secret generic basic-auth \
  --from-literal=basic-auth-user=admin \
  --from-literal=basic-auth-password="$PASSWORD"
```

## 4. Deploy OpenFaaS Components

Apply the OpenFaaS core manifests:

```bash
kubectl apply -f ./yaml/
```

This will deploy:
- OpenFaaS gateway
- faas-netes controller
- Queue worker
- Prometheus & AlertManager (optional)
- Other required system components

## 5. Verify the Deployment

Check whether the OpenFaaS pods are running:

```bash
kubectl get pods -n openfaas
```

All pods should show the status `Running` or `Completed`.


## 6. Access the Gateway

For faas-cli to communicate with gateway:

```bash
export OPENFAAS_URL=http://127.0.0.1:31112
```

Log in via CLI:
```bash
echo $PASSWORD | faas-cli login --username admin --password-stdin
```

## 7. Test the CLI

List all deployed functions (should be empty initially):

```bash
faas-cli list
```
