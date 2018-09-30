NAME  ?=dnsconfig-injector
REGISTRY ?=karampok
ARGS ?=

pack:
	docker build --no-cache -t $(REGISTRY)/$(NAME):v1 .

upload: pack
	docker push $(REGISTRY)/$(NAME):v1
