FROM golang:alpine AS builder
ARG CGO_ENABLED=0
WORKDIR /sources
COPY . /sources
RUN go build -ldflags "-s -w" -o sonar-badge-proxy
RUN apk add -U --no-cache ca-certificates

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /sources/sonar-badge-proxy /opt/sonar-badge-proxy/start
COPY favicon.ico /opt/sonar-badge-proxy/
ENTRYPOINT ["/opt/sonar-badge-proxy/start"]
