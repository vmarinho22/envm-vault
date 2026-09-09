# Multi-stage OCI image: Vue SPA + Go API in one process.
# Compatible with Docker and Podman.

# syntax=docker/dockerfile:1

FROM node:24-alpine AS frontend
WORKDIR /src
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
COPY .deploy/frontend/ /tmp/frontend-env/
RUN cp /tmp/frontend-env/.env.example ./.env && \
    if [ -f /tmp/frontend-env/.env ]; then cp /tmp/frontend-env/.env ./.env; fi
RUN npm run build-only

FROM golang:1.25-alpine AS api-build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY api/go.mod api/go.sum ./
RUN go mod download
COPY api/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=api-build /out/server /app/server
COPY --from=frontend /src/dist /app/ui
ENV FRONTEND_DIST=/app/ui
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/server"]
