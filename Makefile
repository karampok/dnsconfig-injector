NAME  ?=dnsconfig-injector
REGISTRY ?=karampok
ARGS ?=

build:
	CGO_ENABLED=0 GOOS=linux go build -o mutating-dns-webhook-server .

pack: build
	docker build --no-cache -t $(REGISTRY)/$(NAME):latest .

upload: pack
	docker push $(REGISTRY)/$(NAME):v1
