# Fetch
FROM golang:latest AS fetch-stage
COPY go.mod go.sum /app/
WORKDIR /app
RUN go mod download

# Generate
FROM ghcr.io/a-h/templ:latest AS generate-stage
COPY --chown=65532:65532 . /app
WORKDIR /app
RUN ["templ", "generate"]

# Build
FROM golang:latest AS build-stage
COPY --from=generate-stage /app /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/app

# Verify that static files exist in the builder stage
RUN ls -la /app/style

# Deploy
FROM gcr.io/distroless/base-debian12 AS deploy-stage
WORKDIR /
COPY --from=build-stage /app/app /app
COPY --from=build-stage /app/style /style
COPY --from=build-stage /app/assets /assets
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app"]
