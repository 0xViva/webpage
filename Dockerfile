FROM golang:1.24 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o /server main.go

# Verify that static files exist in the builder stage
RUN ls -la /app/style

FROM gcr.io/distroless/base-debian11 AS final

COPY --from=builder /server /server

# Copy static files to the final image
COPY --from=builder /app/style /style
COPY --from=builder /app/assets /assets


ENV PORT=8080
EXPOSE $PORT

ENTRYPOINT ["/server"]
