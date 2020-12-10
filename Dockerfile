FROM golang:alpine AS build

WORKDIR /source

COPY . /source

RUN  go build

FROM alpine 

COPY --from=build /source/sonar-badge-proxy /sonar-badge-proxy

ENTRYPOINT [ "/sonar-badge-proxy" ]

