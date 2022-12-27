# syntax=docker/dockerfile:1
FROM golang:1.19.4-bullseye AS Builder
WORKDIR /

COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:3.14
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app ./
COPY config.yaml ./

CMD ["./app"]