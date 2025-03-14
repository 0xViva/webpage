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
FROM golang:1.24 AS build-stage
WORKDIR /app
COPY --from=fetch-stage /app /app 
COPY --from=generate-stage /app /app  
RUN CGO_ENABLED=0 GOOS=linux go build -o /server main.go

# Deploy minimal image
FROM gcr.io/distroless/base-debian11 AS deploy-stage
COPY --from=build-stage /server /server
# include all static files since they're not embedded in the binary
COPY --from=build-stage /app/style /style
COPY --from=build-stage /app/assets /assets
ENV PORT=8080
EXPOSE $PORT
ENTRYPOINT ["/server"]

