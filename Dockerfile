# build stage
# golang:1.19.4-bullseye
FROM golang@sha256:04f76f956e51797a44847e066bde1341c01e09054d3878ae88c7f77f09897c4d AS build-env

#COPY certs/ /usr/local/share/ca-certificates/
#RUN update-ca-certificates

RUN mkdir /build
COPY . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o dnsconfig-injector .

# final stage
FROM scratch
WORKDIR /app
COPY --from=build-env /build/dnsconfig-injector /app/

CMD ["/app/dnsconfig-injector"]
