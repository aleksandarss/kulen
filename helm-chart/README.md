# Recipe App Helm Deployment

This chart deploys the Recipe App (frontend, backend, and optionally PostgreSQL) on Kubernetes.

---

## ✅ Prerequisites

- Kubernetes cluster (with Ingress controller like NGINX installed)
- Helm installed
- Optional: TLS certificate and key for HTTPS

---

## 🛠 Manual Steps

### 1. Create TLS Secret (if TLS is enabled)

If you want to use HTTPS for your Ingress, create a TLS secret in your cluster:

```bash
kubectl create secret tls recipe-tls \
  --cert=path/to/tls.crt \
  --key=path/to/tls.key
```

```yaml
ingress:
  enabled: true
  host: your.domain.com
  tlsSecret: recipe-tls
```

### 2. Persistent Volume for PostgreSQL

Kubernetes will automatically create a PersistentVolumeClaim (PVC) if postgresql.enabled: true.

Ensure your cluster has a default StorageClass or set one explicitly:
```yaml
postgresql:
  enabled: true
  storageClass: standard
  storage: 1Gi
```
No need to manually create the volume unless your cluster doesn't support dynamic provisioning.

### 3. Configure External PostgreSQL (optional)
If you disable in-cluster PostgreSQL (postgresql.enabled: false), configure connection settings:

```yaml
backend:
  env:
    DB_HOST: your-external-db-host
    DB_PORT: 5432
    DB_USER: your-db-user
    DB_PASSWORD: your-db-password
    DB_NAME: your-db-name
```

### 4. Deploy
```bash
helm upgrade recipe-app ./chart -f values.yaml
```

### 5. Upgrade
```bash
helm upgrade recipe-app ./chart -f values.yaml
```

### Notes

Make sure Docker images are pushed to a registry accessible by your cluster (e.g., GitHub Container Registry or Docker Hub).

Ingress assumes a running Ingress controller in the cluster.

For HTTPS, DNS should resolve to your cluster and TLS secret must exist.

