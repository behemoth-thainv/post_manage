# Frontend Deployment Guide

## Tổng quan

Frontend được build bằng Vue.js + Vite và được triển khai lên Kubernetes sử dụng Nginx để serve static files.

## Các file liên quan

### Frontend Files
- `frontend/Dockerfile` - Multi-stage Docker build
- `frontend/nginx.conf` - Nginx configuration cho SPA routing

### Kubernetes Manifests
- `k8s/frontend-deployment.yaml` - Deployment và Service cho frontend
- `k8s/ingress.yaml` - Ingress routing (đã được cập nhật)

## Hướng dẫn triển khai

### 1. Build Docker Image

```bash
cd frontend

# Build image
docker build -t asia-northeast1-docker.pkg.dev/glossy-attic-443507-u6/docker/frontend:latest .

# Push to registry
docker push asia-northeast1-docker.pkg.dev/glossy-attic-443507-u6/docker/frontend:latest
```

### 2. Deploy lên Kubernetes

```bash
cd ../k8s

# Deploy frontend
kubectl apply -f frontend-deployment.yaml

# Verify deployment
kubectl get pods -n demo
kubectl get svc -n demo
```

### 3. Update Ingress

Ingress đã được cấu hình để:
- Route `/api/*` → backend-rails service
- Route `/*` (catch-all) → frontend service

```bash
# Apply updated ingress
kubectl apply -f ingress.yaml

# Verify ingress
kubectl get ingress -n demo
kubectl describe ingress demo-ingress -n demo
```

## Kiến trúc

### Multi-stage Build

```
Stage 1 (builder): Node.js 18 Alpine
  ├── npm ci - Install dependencies
  ├── npm run build - Build production bundle
  └── Output: /app/dist

Stage 2 (production): Nginx Alpine
  ├── Copy dist from builder
  ├── Copy nginx.conf
  └── Serve on port 80
```

### Nginx Configuration

- **SPA Routing**: Redirect tất cả requests về `index.html`
- **Static Asset Caching**: Cache 1 năm cho JS, CSS, images
- **Gzip Compression**: Giảm kích thước response
- **Health Check**: Endpoint `/health` cho Kubernetes probes

### Kubernetes Resources

**Deployment:**
- 2 replicas để đảm bảo high availability
- Resource limits: 200m CPU, 256Mi memory
- Resource requests: 50m CPU, 64Mi memory
- Readiness/Liveness probes trên `/health`

**Service:**
- Type: ClusterIP (internal only)
- Port: 80

**Ingress:**
- Path-based routing
- GCE ingress class (Google Cloud Load Balancer)

## Troubleshooting

### Check Pod Status
```bash
kubectl get pods -n demo -l app=frontend
kubectl logs -n demo -l app=frontend
```

### Check Service
```bash
kubectl get svc -n demo frontend
kubectl describe svc -n demo frontend
```

### Check Ingress
```bash
kubectl get ingress -n demo
kubectl describe ingress -n demo demo-ingress
```

### Test Locally
```bash
# Build and run locally
docker build -t frontend-test .
docker run -p 8080:80 frontend-test

# Access at http://localhost:8080
```

## Environment Variables

Frontend không cần environment variables riêng vì API endpoint được proxy qua Ingress. Tất cả requests đến `/api/*` sẽ tự động được route đến backend.

## Scaling

Để scale frontend:

```bash
kubectl scale deployment frontend -n demo --replicas=5
```

## Updates

Để update frontend:

```bash
# Build new image with tag
docker build -t asia-northeast1-docker.pkg.dev/glossy-attic-443507-u6/docker/frontend:v1.0.1 .
docker push asia-northeast1-docker.pkg.dev/glossy-attic-443507-u6/docker/frontend:v1.0.1

# Update deployment
kubectl set image deployment/frontend frontend=asia-northeast1-docker.pkg.dev/glossy-attic-443507-u6/docker/frontend:v1.0.1 -n demo

# Or rollout restart
kubectl rollout restart deployment/frontend -n demo
kubectl rollout status deployment/frontend -n demo
```
