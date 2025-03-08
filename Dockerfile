FROM golang:1.24 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o /server main.go

RUN ls

FROM gcr.io/distroless/base-debian11 AS final

COPY --from=builder /server /server

ENV PORT=8080
EXPOSE $PORT

ENTRYPOINT ["/server"]