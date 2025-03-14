# Fetch dependencies
FROM golang:1.24 AS fetch-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Generate
FROM ghcr.io/a-h/templ:latest AS generate-stage
COPY --chown=65532:65532 . /app
WORKDIR /app
RUN ["templ", "generate"]

# Build the server binary
FROM golang:1.24 AS builder
WORKDIR /app
COPY --from=fetch-stage /app /app 
COPY --from=generate-stage /app /app  
RUN CGO_ENABLED=0 GOOS=linux go build -o /server main.go

# Verify static files exist
RUN ls -la /app/style

# Deploy minimal image
FROM gcr.io/distroless/base-debian11 AS final
COPY --from=builder /server /server
COPY --from=builder /app/style /style
COPY --from=builder /app/assets /assets

ENV PORT=8080
EXPOSE $PORT

ENTRYPOINT ["/server"]

