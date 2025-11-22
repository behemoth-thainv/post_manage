# Training Outline: Deploy Web App to GKE (Google Kubernetes Engine)

Tài liệu này cung cấp giàn ý cho buổi training về việc đưa ứng dụng Web (Frontend + Backend) lên GKE sử dụng các dịch vụ của Google Cloud Platform (GCP).

## 1. Giới thiệu & Kiến trúc (Introduction & Architecture)
- **Mục tiêu**: Hiểu rõ quy trình CI/CD và cách vận hành ứng dụng trên Kubernetes.
- **Luồng hoạt động (Workflow)**:
  1. **Code**: Developer push code lên Git.
  2. **CI/CD**: Cloud Build tự động build Docker Image.
  3. **Registry**: Docker Image được lưu trữ tại Artifact Registry.
  4. **Deploy**: Cloud Build hoặc Developer deploy manifest lên GKE.
  5. **Runtime**: Ứng dụng chạy trên Pods, kết nối Database (Cloud SQL) và Cache (Redis).

## 2. Chuẩn bị (Prerequisites)
- **Công cụ cần thiết**:
  - Google Cloud SDK (`gcloud`).
  - Kubernetes CLI (`kubectl`).
  - Docker (để test local).
- **Quyền hạn (IAM)**:
  - Kubernetes Engine Admin.
  - Cloud Build Editor.
  - Artifact Registry Administrator.
  - Service Account User.

## 3. Các dịch vụ GCP (GCP Services)

### 3.1. Artifact Registry
- **Vai trò**: Lưu trữ Docker Images (thay thế Container Registry cũ).
- **Thực hành**:
  - Tạo Repository (Format: Docker).
  - Cấu hình Docker auth với gcloud.
  - Build và Push image thủ công lên Registry.

### 3.2. Cloud Build (CI/CD)
- **Vai trò**: Tự động hóa quy trình Build & Deploy.
- **Thực hành**:
  - Tạo file `cloudbuild.yaml`.
  - Cấu hình Steps: Build Image -> Push Image -> Update K8s Manifest -> Deploy.
  - Thiết lập Trigger: Tự động chạy khi có commit vào nhánh `main` hoặc `develop`.

### 3.3. Cloud SQL (Database)
- **Vai trò**: Dịch vụ Database quản lý (PostgreSQL/MySQL).
- **Thực hành**:
  - Tạo Instance Cloud SQL.
  - Tạo Database và User.
  - **Kết nối từ GKE**:
    - Cách 1: Sử dụng Cloud SQL Auth Proxy (Sidecar container).
    - Cách 2: Sử dụng Private IP (VPC Native Cluster).

### 3.4. Cloud Memorystore (Redis)
- **Vai trò**: Dịch vụ Redis quản lý cho caching/session.
- **Thực hành**:
  - Tạo Instance Redis.
  - Lấy địa chỉ IP và Port.
  - Cấu hình ứng dụng để kết nối Redis.

## 4. Kubernetes (GKE) Concepts

### 4.1. ConfigMap & Secret
- **ConfigMap**: Lưu cấu hình không nhạy cảm (Environment variables, config files).
  - Ví dụ: `DB_HOST`, `API_URL`, `DEBUG_MODE`.
- **Secret**: Lưu thông tin nhạy cảm (được mã hóa).
  - Ví dụ: `DB_PASSWORD`, `REDIS_PASSWORD`, `API_KEYS`.
- **Thực hành**: Tạo ConfigMap/Secret từ file hoặc command line và mount vào Pod.

### 4.2. Workloads (Deployment)
- **Deployment**: Quản lý số lượng bản sao (replicas) của ứng dụng (Pods).
- **Thực hành**:
  - Viết file `deployment.yaml` cho Frontend và Backend.
  - Cấu hình `image` (lấy từ Artifact Registry).
  - Cấu hình `resources` (CPU/Memory limits).
  - Cấu hình `env` (lấy từ ConfigMap/Secret).

### 4.3. Services & Networking
- **Service**: Expose ứng dụng để các thành phần khác truy cập.
  - `ClusterIP`: Chỉ truy cập nội bộ (Backend, Redis).
  - `LoadBalancer`: Expose ra Internet (Frontend).
- **Ingress** (Nâng cao): Quản lý routing HTTP/HTTPS, SSL termination.

## 5. K8s Concepts (Advanced)

### 5.1. Auto-scaling (HPA)
- **Horizontal Pod Autoscaler**: Tự động tăng/giảm số lượng Pod dựa trên CPU/Memory usage.
- **Thực hành**:
  - Cấu hình `resources.requests` và `limits` cho Pod.
  - Tạo HPA rule: `kubectl autoscale deployment [NAME] --cpu-percent=50 --min=1 --max=10`.
  - Stress test để xem Pod tự động scale up.

### 5.2. Security Best Practices
- **Workload Identity**: Cho phép Pod truy cập GCP Services (Cloud SQL, Storage) an toàn mà **không cần Service Account Key (JSON)**.
  - **Cơ chế**: Liên kết (bind) Kubernetes Service Account (KSA) với Google Service Account (GSA). Pod sẽ dùng token ngắn hạn để xác thực.
  - **Lợi ích**: Không lo lộ key, không cần rotate key thủ công, tuân thủ nguyên tắc Least Privilege.
  - **Cấu hình mẫu**:
    ```bash
    # 1. Tạo GSA
    gcloud iam service-accounts create k8s-demo-sa \
        --project=glossy-attic-443507-u6

    # 2. Cấp quyền cho GSA (DB, Logs, Storage)
    # Cloud SQL Client
    gcloud projects add-iam-policy-binding glossy-attic-443507-u6 \
        --member "serviceAccount:k8s-demo-sa@glossy-attic-443507-u6.iam.gserviceaccount.com" \
        --role "roles/cloudsql.client"

    # Cloud Logging Writer
    gcloud projects add-iam-policy-binding glossy-attic-443507-u6 \
        --member "serviceAccount:k8s-demo-sa@glossy-attic-443507-u6.iam.gserviceaccount.com" \
        --role "roles/logging.logWriter"

    # Storage Object User (S3)
    gcloud projects add-iam-policy-binding glossy-attic-443507-u6 \
        --member "serviceAccount:k8s-demo-sa@glossy-attic-443507-u6.iam.gserviceaccount.com" \
        --role "roles/storage.objectUser"

    # 3. Tạo KSA (Kubernetes Service Account)
    kubectl create serviceaccount post-manage-ksa --namespace demo

    # 4. Liên kết GSA với KSA (Binding)
    # Giả sử K8s Namespace là 'demo' và KSA là 'post-manage-ksa'
    gcloud iam service-accounts add-iam-policy-binding k8s-demo-sa@glossy-attic-443507-u6.iam.gserviceaccount.com \
        --role roles/iam.workloadIdentityUser \
        --member "serviceAccount:glossy-attic-443507-u6.svc.id.goog[demo/post-manage-ksa]"

    # 5. Annotate KSA
    kubectl annotate serviceaccount post-manage-ksa \
        iam.gke.io/gcp-service-account=k8s-demo-sa@glossy-attic-443507-u6.iam.gserviceaccount.com
    ```
- **Network Policies**: Kiểm soát traffic giữa các Pods (ví dụ: chỉ cho phép Backend gọi Database).

## 6. Observability (Giám sát & Vận hành)

### 6.1. Cloud Logging
- **Logs Explorer**: Xem log tập trung của tất cả các Pods.
- **Thực hành**:
  - Filter log theo `resource.labels.cluster_name` và `labels.k8s-pod/app`.
  - Cách viết log chuẩn (JSON format) để dễ query.

### 6.2. Cloud Monitoring
- **Metrics**: Theo dõi CPU, Memory, Disk I/O của Cluster và Pods.
- **Alerting**: Cảnh báo khi Pod crash hoặc Resource usage quá cao.

## 7. Troubleshooting (Gỡ lỗi thường gặp)
- **CrashLoopBackOff**: Pod khởi động nhưng bị crash ngay lập tức (thường do lỗi code hoặc thiếu config).
- **ImagePullBackOff**: Không thể kéo Docker image (sai tên image, sai tag, hoặc chưa cấp quyền Artifact Registry).
- **Pending**: Không đủ tài nguyên (CPU/RAM) trên Node để chạy Pod.
- **Debug Command**:
  - `kubectl describe pod [POD_NAME]`: Xem events chi tiết.
  - `kubectl logs [POD_NAME] -f`: Xem log realtime.
  - `kubectl exec -it [POD_NAME] -- sh`: Chui vào container để debug.

## 8. Quy trình thực hành (Hands-on Flow)

1.  **Setup Environment**: Cài đặt `gcloud`, `kubectl`, tạo Project.
2.  **Containerize App**: Viết `Dockerfile` cho FE và BE.
3.  **Manual Push**: Build và push image lên Artifact Registry thủ công để kiểm tra.
4.  **Setup Infra**: Tạo Cloud SQL và Redis instance.
5.  **K8s Manifests**:
    - Tạo `configmap.yaml` chứa biến môi trường chung.
    - Tạo `secret.yaml` (tên `app-secrets`) chứa DB credentials và Rails secrets.
    - Tạo `backend-go-deployment.yaml` và `backend-rails-deployment.yaml`.
    - Tạo `ingress.yaml` để điều hướng traffic.
6.  **Deploy**:
    ```bash
    kubectl apply -f k8s/configmap.yaml
    kubectl apply -f k8s/secret.yaml
    kubectl apply -f k8s/backend-go-deployment.yaml
    kubectl apply -f k8s/backend-rails-deployment.yaml
    kubectl apply -f k8s/ingress.yaml
    ```
7.  **Verify**: Kiểm tra logs (`kubectl logs`), kiểm tra kết nối DB/Redis.
8.  **Automation**: Cấu hình Cloud Build để tự động hóa quy trình trên.

---

## 9. Reference Commands (Các lệnh thường dùng)

### Auth
```bash
gcloud auth login
gcloud config set project [PROJECT_ID]
```

### Fetch GKE credentials
```bash
gcloud container clusters get-credentials [CLUSTER_NAME] --zone [COMPUTE_ZONE]
# or
gcloud container clusters get-credentials [CLUSTER_NAME] --region [COMPUTE_REGION]
```

### Create a namespace
```bash
kubectl create namespace [NAMESPACE]
```

### Deploy to GKE
```bash
kubectl apply -f [DEPLOYMENT_YAML]
```

### Expose the service
```bash
kubectl expose deployment [DEPLOYMENT_NAME] --type=LoadBalancer --port=[PORT]
```

### Get the external IP
```bash
kubectl get svc [SERVICE_NAME]
```

### Delete resources
```bash
kubectl delete svc [SERVICE_NAME]
kubectl delete deployment [DEPLOYMENT_NAME]
kubectl delete namespace [NAMESPACE]
```

### Delete the GKE cluster
```bash
gcloud container clusters delete [CLUSTER_NAME] --zone [COMPUTE_ZONE]
```
