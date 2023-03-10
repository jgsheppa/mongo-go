# syntax=docker/dockerfile:1
FROM golang:1.19.4-alpine3.17 AS builder
WORKDIR /

COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM golang:1.19.4-alpine3.17
RUN apk --no-cache add ca-certificates
RUN --mount=type=secret,id=config_yaml,dst=/etc/secrets/config.yaml cat /etc/secrets/config.yaml
WORKDIR /root/
COPY --from=builder /app ./


CMD ["./app"]