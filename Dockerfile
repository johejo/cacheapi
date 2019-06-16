FROM golang:1.12-alpine AS builder
LABEL maintainer="mitsuo_h@outlook.com"

ENV GO111MODULE=on \
    GOOS=linux
RUN apk add --no-cache git make upx
WORKDIR /go/src/github.com/johejo/cacheapi
COPY . .
RUN make min-binary

FROM alpine:3.9
COPY --from=builder /go/src/github.com/johejo/cacheapi/out/cacheapi /cacheapi
EXPOSE 8888
ENTRYPOINT [ "/cacheapi" ]
