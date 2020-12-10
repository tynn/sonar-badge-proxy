FROM golang:alpine AS build

WORKDIR /source

RUN apk add git

RUN git clone https://github.com/blink38/sonar-badge-proxy

RUN cd sonar-badge-proxy && go build



FROM alpine 

LABEL MAINTAINER=<matthieu.marc@gmail.com>

COPY --from=build /source/sonar-badge-proxy/sonar-badge-proxy /sonar-badge-proxy

ENTRYPOINT [ "/sonar-badge-proxy" ]

