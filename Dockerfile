# Build the React frontend
FROM node:18 AS frontend-build

# Set working directory for the frontend
WORKDIR /frontend

# Install dependencies and build the React app
COPY frontend/package.json frontend/package-lock.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# Build the Go backend with custom version
FROM debian:bullseye-slim AS backend-build

WORKDIR /backend

# Install dependencies
RUN apt-get update && apt-get install -y \
    wget \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Download and install Go 1.23.4
RUN wget https://go.dev/dl/go1.23.4.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz \
    && rm go1.23.4.linux-amd64.tar.gz

# Add Go to PATH
ENV PATH="/usr/local/go/bin:$PATH"

# Install dependencies and build the Go backend
COPY backend/go.mod backend/go.sum ./
RUN go mod tidy
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/backend .

# Combine the frontend and backend
FROM debian:bullseye-slim

# Set working directory for the final image
WORKDIR /app

# Install necessary dependencies (e.g., for serving the frontend)
RUN apt-get update && apt-get install -y \
    ca-certificates \
    python3 \
    && rm -rf /var/lib/apt/lists/*

# Copy the Go backend from the build stage
COPY --from=backend-build /app/backend /app/backend

# Copy the React frontend build from the build stage
COPY --from=frontend-build /frontend/dist /app/frontend

# Expose ports for the frontend and backend
EXPOSE 3000 8080

# Start the backend and serve the frontend
CMD ["sh", "-c", "cd /app/frontend && python3 -m http.server 3000 & /app/backend"]
