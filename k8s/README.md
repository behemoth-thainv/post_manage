**Overview**:

This folder contains Kubernetes manifests for deploying the sample_crud application to GKE using Cloud SQL for MySQL.

**Files of interest**:
- `namespace.yaml` - Kubernetes namespace `sample-crud` used for all resources.
- `cloudsql-configmap.yaml` - Non-sensitive database config (host/port, environment).
- `cloudsql-secret.yaml` - Database credentials and `INSTANCE_CONNECTION_NAME` (do not commit real secrets).
- `cloudsql-sa-secret.yaml` - (Optional) Service account JSON for the Cloud SQL Proxy. Prefer Workload Identity instead of placing JSON in the repo.
- `backend-deployment.yaml` - Rails backend Deployment. Uses a Cloud SQL Proxy sidecar and reads settings via `envFrom`.
- `backend-go-deployment.yaml` - Go backend Deployment. Uses a Cloud SQL Proxy sidecar and reads settings via `envFrom`.
- `frontend-deployment.yaml` and `frontend-service.yaml` - Frontend resources.
- `ingress.yaml` - Ingress configured for GCE; update host and optionally add ManagedCertificate.
- `kustomization.yaml` - Aggregates manifests; apply with `kubectl apply -k k8s/`.

**Important notes & recommended workflow**:

1. Replace placeholders
   - Update `gcr.io/PROJECT_ID/...` image references in the deployment files to your actual image locations (Container Registry or Artifact Registry).
   - Update `cloudsql-secret.yaml`'s `INSTANCE_CONNECTION_NAME` to your Cloud SQL instance connection name: `PROJECT:REGION:INSTANCE`.
   - Update the Cloud SQL proxy args in deployments: replace `PROJECT_ID:REGION:INSTANCE_NAME` with your instance connection name.

2. Secrets and credentials
   - Do NOT commit service-account keys or real passwords into source control. Create the secrets directly in the cluster:

```bash
kubectl create secret generic cloudsql-secret \
  --from-literal=DB_NAME=sample_crud_db \
  --from-literal=DB_USER=sample_crud_user \
  --from-literal=DB_PASSWORD='YOUR_PASSWORD' \
  --from-literal=INSTANCE_CONNECTION_NAME='PROJECT:REGION:INSTANCE' \
  -n sample-crud

# If you are not using Workload Identity and need to use a service account key file:
kubectl create secret generic cloudsql-sa --from-file=credentials.json=/path/to/key.json -n sample-crud
```

   - Recommended: enable GKE Workload Identity and grant the GKE service account access to Cloud SQL. This avoids embedding JSON keys.

3. Build & push images
   - Rails backend: build using `backend/Dockerfile` then push to your registry.
   - Go backend: build using `backend_go/Dockerfile` then push.
   - Frontend: build static assets and push an image that serves them (nginx or static file server).

4. Apply manifests

```bash
# Ensure kubectl points to your GKE cluster
kubectl apply -k k8s/
```

5. TLS and Ingress
   - The `ingress.yaml` is annotated for the GCE ingress controller. To enable HTTPS, reserve a static IP and create a ManagedCertificate, then annotate the Ingress accordingly.

6. Postgres manifests
   - Postgres manifests previously present in this directory are deprecated (Cloud SQL is used). They are left as marked/deprecated files; you can safely delete them if you prefer.

**Security reminder**: Move DB passwords and service-account JSON out of git. Prefer Workload Identity + Cloud SQL IAM-based authentication.

If you want, I can:
- Replace placeholders with your GCP project/instance details and update the files.
- Add example `Dockerfile` build-and-push commands for each service.
- Configure a ManagedCertificate + static IP for HTTPS in `ingress.yaml`.
