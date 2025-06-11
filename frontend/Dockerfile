# ---- Build step ----
FROM node:18 AS builder
WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm install

# Copy source and build
COPY . .
RUN npm run build

# ---- Serve step ----
FROM nginx:alpine AS runner
COPY --from=builder /app/dist /usr/share/nginx/html

# Optional: Replace default config
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
