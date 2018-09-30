FROM golang:latest AS builder
COPY . $GOPATH/src/mypackage/myapp/
WORKDIR $GOPATH/src/mypackage/myapp/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /mutating-dns-webhook-server .

FROM alpine:latest
COPY --from=builder /mutating-dns-webhook-server mutating-dns-webhook-server
ENTRYPOINT ["./mutating-dns-webhook-server"]
