SHELL=/bin/bash

IMAGE_PREFIX=hipster
DOCKER_SERVER=gcr.io/metal-cascade-146909

VERSION = $(version)

PROXY_BINARY_LINUX=proxy/proxy_linux
PROXY_BINARY=proxy/proxy
API_BINARY_LINUX=api/api_linux
API_BINARY=api/api

LDFLAGS=-ldflags "-X mainersion=${version}"

.DEFAULT_GOAL: all

$(API_BINARY): api/main.go
	go build -o $@ api/main.go

$(PROXY_BINARY): proxy/main.go
	go build -o $@ proxy/main.go


$(API_BINARY_LINUX): api/main.go
ifeq ($(VERSION),)
	$(error VERSION is not set)
else
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o $@ api/main.go
endif

$(PROXY_BINARY_LINUX): proxy/main.go
ifeq ($(VERSION),)
	$(error VERSION is not set)
else
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o $@ proxy/main.go
endif

.PHONY: build_api_app
build_api_app: $(API_BINARY_LINUX)

#.PHONY: build_api
#build_api: $(API_BINARY)

.PHONY: build_proxy_app
build_proxy_app: $(PROXY_BINARY_LINUX)

.PHONY: build_proxy
build_proxy: $(PROXY_BINARY)

.PHONY: build_api_image
build_api_image:
ifeq ($(VERSION),)
	$(error VERSION is not set)
else
	cd api;\
	docker build -t $(IMAGE_PREFIX)/api:$(VERSION) .; \
	docker tag $(IMAGE_PREFIX)/api:$(VERSION) $(DOCKER_SERVER)/api:$(VERSION); \
	gcloud docker -- push $(DOCKER_SERVER)/api:$(VERSION); \
	cd -;
	rm $(API_BINARY_LINUX);
endif

.PHONY: build_proxy_image
build_proxy_image:
ifeq ($(VERSION),)
	$(error VERSION is not set)
else
	cd proxy;\
	docker build -t $(IMAGE_PREFIX)/proxy:$(VERSION) .; \
	docker tag $(IMAGE_PREFIX)/proxy:$(VERSION) $(DOCKER_SERVER)/proxy:$(VERSION); \
	gcloud docker -- push $(DOCKER_SERVER)/proxy:$(VERSION); \
	cd -;
	rm $(PROXY_BINARY_LINUX);
endif


.PHONY: build_api
build_api: build_api_app build_api_image

.PHONY: build_proxy
build_proxy: build_proxy_app build_proxy_image


.PHONY: build
build: $(PROXY_BINARY) $(API_BINARY)

.PHONY: deps
deps: glide.lock glide.yaml

glide.lock: glide.yaml
	glide i
	touch glide.lock
