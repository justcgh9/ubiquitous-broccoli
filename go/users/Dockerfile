FROM golang:1.24 AS builder
WORKDIR /app

ARG SERVER_PORT
ARG CONFIG_PATH

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/
COPY ${CONFIG_PATH} config/prod.yml

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/users ./cmd/users

FROM alpine:latest
WORKDIR /app

ARG SERVER_PORT

COPY --from=builder /app/bin/users ./bin/users
COPY --from=builder /app/config/prod.yml ./config/prod.yml

EXPOSE ${SERVER_PORT}

CMD ["./bin/users", "--config=./config/prod.yml"]