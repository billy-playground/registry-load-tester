FROM docker.io/library/golang:1.23.4-alpine AS builder
RUN apk add make
ADD . /src
RUN rm -rf /src/bin
WORKDIR /src
RUN make build  # Static build for amd64
RUN mv /src/bin/rlt /go/bin/rlt
RUN mv /src/bin/assets /go/bin/assets

FROM alpine:3.17.1
RUN apk add file
RUN apk --update add ca-certificates file

RUN mkdir /app
COPY --from=builder /go/bin/rlt /app/rlt
RUN chmod +x /app/rlt
COPY --from=builder /go/bin/assets /app/assets

WORKDIR /app
ENTRYPOINT ["/app/rlt"]