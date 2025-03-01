FROM golang:1.24 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /server main.go

FROM gcr.io/distroless/base-debian11 as final

COPY --from=builder /server /server

ENV PORT 8080
EXPOSE $PORT

ENTRYPOINT ["/server"]