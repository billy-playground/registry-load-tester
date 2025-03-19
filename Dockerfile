FROM docker.io/library/golang:1.23.4-alpine AS builder
RUN apk add make
ADD . /src
WORKDIR /src
RUN make build  # Static build for amd64
RUN mv /src/bin/test /go/bin/test
RUN mv /src/bin/assets /go/bin/assets

FROM alpine:3.17.1
RUN apk add file
RUN apk --update add ca-certificates file

RUN mkdir /app
COPY --from=builder /go/bin/test /app/test
RUN chmod +x /app/test
COPY --from=builder /go/bin/assets /app/assets

WORKDIR /app
ENTRYPOINT ["/app/test"]