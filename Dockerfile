# syntax=docker/dockerfile:1
FROM 1.19.4-alpine3.17 AS builder
WORKDIR /

COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM 1.19.4-alpine3.17
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app ./
COPY config.yaml ./

# create group and user
RUN groupadd -r gopher && useradd -g gopher gopher
# set ownership and permissions
RUN chown -R gopher:gopher /app
# switch user
USER gopher

CMD ["./app"]