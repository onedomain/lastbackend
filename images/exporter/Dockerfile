# Build manifest
FROM golang:1.11-alpine3.8 as build

RUN apk add --no-cache ca-certificates

RUN apk add --no-cache iptables \
    linux-headers \
    gcc \
    musl-dev \
    iproute2

RUN set -ex \
	&& apk add --no-cache --virtual .build-deps \
    bash \
    git  \
    make \
	\
	&& rm -rf /*.patch

ENV GOPATH /go
ENV GOROOT /usr/local/go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

ADD . /go/src/github.com/onedomain/lastbackend
WORKDIR /go/src/github.com/onedomain/lastbackend

RUN make APP=exporter build && make APP=exporter install
RUN apk del --purge .build-deps

WORKDIR /go/bin
RUN rm -rf /go/pkg \
    && rm -rf /go/src \
    && rm -rf /var/cache/apk/*


# Production manifest
FROM alpine:3.8 as production

RUN apk add --no-cache ca-certificates \
  iptables \
  iproute2

COPY --from=build /usr/bin/exporter /usr/bin/exporter

EXPOSE 2963