# ---- Build step ----
FROM node:18 AS builder
WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

ARG VITE_API_URL_ARG="https://kulen.mithrandir.calic.cloud/api"
ENV VITE_API_URL=${VITE_API_URL_ARG}

RUN npm run build

# ---- Serve step ----
FROM nginx:alpine AS runner
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
