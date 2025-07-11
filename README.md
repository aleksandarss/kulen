# Kulen

This repository contains a small recipe application split into a Go backend and a Vue 3 frontend. The project can be started locally with Docker Compose or deployed to Kubernetes using the provided Helm chart.

## Coverage

![Backend coverage](https://img.shields.io/badge/backend-80%25-brightgreen)
![Frontend coverage](https://img.shields.io/badge/frontend-100%25-brightgreen)

Coverage percentages are computed in the CI run and are expected to stay above 80%.

---

## Local Development

### Prerequisites

- Go 1.24
- Node 18
- Docker & Docker Compose

### Backend

The API server lives in `backend/`. It reads configuration from environment variables:

- `DATABASE_URL` – PostgreSQL connection string
- `PORT` – HTTP port (defaults to `8080`)
- `JWT_SECRET` – secret used for signing JWT tokens

Run the backend directly with Go:

```bash
cd backend
DATABASE_URL=postgres://user:password@localhost:5432/recipeapp \
JWT_SECRET=supersecret go run ./cmd/server
```

### Frontend

The web client is in `frontend/` and uses Vite. The `VITE_API_URL` variable must point to the backend API:

```bash
cd frontend
npm install
VITE_API_URL=http://localhost:8080/api npm run dev
```

### Using docker-compose

Running everything at once can be done with Docker Compose. This spins up PostgreSQL, the backend and the frontend:

```bash
docker-compose up --build
```

The application will then be reachable at `http://localhost:8081` while the API is served on `http://localhost:8080`.

### Running tests

- Backend: `go test ./...`
- Frontend: `npm test`

---

## Docker Images

Both services contain Dockerfiles. To build them separately:

```bash
# Backend
cd backend
docker build -t kulen-backend .

# Frontend
cd ../frontend
docker build -t kulen-frontend .
```

---

## Helm Chart

The `helm-chart/` directory contains everything required to deploy the application on Kubernetes.

### Chart Structure

- `Chart.yaml` – basic chart metadata
- `values.yaml` – default values for images, environment variables, ingress and optional PostgreSQL
- `templates/` – Kubernetes manifests for the frontend, backend, PostgreSQL, services and ingress

The chart deploys two Deployments (frontend and backend), Services for both, an optional PostgreSQL StatefulSet/Service and an Ingress resource that routes `/api` to the backend and the root path to the frontend.

### Deploying

Adjust `values.yaml` if required (e.g. image tags, ingress host or disabling PostgreSQL) and then run:

```bash
helm upgrade --install recipe-app ./helm-chart -f helm-chart/values.yaml
```

This will create or upgrade the release named `recipe-app`.

For further details see [`helm-chart/README.md`](helm-chart/README.md).

