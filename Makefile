# Copyright Jetstack Ltd. See LICENSE for details.
BINDIR ?= $(PWD)/bin
PATH   := $(BINDIR):$(PATH)

HACK_DIR     ?= hack

REGISTRY := quay.io/jetstack
IMAGE_NAME := vault-helper
IMAGE_TAGS := canary
BUILD_TAG := build

BUILD_IMAGE_NAME := golang:1.9.2

CI_COMMIT_TAG ?= unknown
CI_COMMIT_SHA ?= unknown

VAULT_VERSION := 0.9.5
VAULT_HASH := f6dbc9fdac00598d2a319c9b744b85bf17d9530298f93d29ef2065bc751df099

help:
	# all       - runs verify, build targets
	# test      - runs go_test target
	# build     - runs generate, and then go_build targets
	# generate  - generates mocks and assets files
	# verify    - verifies generated files & scripts
	# image     - build docker image

.PHONY: all test verify

verify: generate go_verify

all: verify build

build: generate go_build

test: go_test

generate: go_generate

go_verify: go_fmt go_vet verify_boilerplate go_test

go_test:
	go test $$(go list ./pkg/... ./cmd/...)

go_fmt:
	@set -e; \
	GO_FMT=$$(git ls-files *.go | grep -v 'vendor/' | xargs gofmt -d); \
	if [ -n "$${GO_FMT}" ] ; then \
		echo "Please run go fmt"; \
		echo "$$GO_FMT"; \
		exit 1; \
	fi

go_vet:
	go vet $$(go list ./pkg/... ./cmd/...)

verify_boilerplate:
	$(HACK_DIR)/verify-boilerplate.sh

go_build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -X main.version=$(CI_COMMIT_TAG) -X main.commit=$(CI_COMMIT_SHA) -X main.date=$(shell date -u +%Y-%m-%d_%H:%M:%S)' -o vault-helper_linux_amd64

image:
	docker build -t $(REGISTRY)/$(IMAGE_NAME):$(BUILD_TAG) .

save:
	docker save $(REGISTRY)/$(IMAGE_NAME):$(BUILD_TAG) -o vault-helper-image.tar

bin/mockgen:
	mkdir -p $(BINDIR)
	go build -o $(BINDIR)/mockgen ./vendor/github.com/golang/mock/mockgen

bin/vault:
	which unzip || apt-get install -y unzip
	mkdir -p $(BINDIR)
	curl -sL  https://releases.hashicorp.com/vault/$(VAULT_VERSION)/vault_$(VAULT_VERSION)_linux_amd64.zip > $(BINDIR)/vault.zip
	echo "$(VAULT_HASH)  $(BINDIR)/vault.zip" | sha256sum  -c
	cd $(BINDIR) && unzip vault.zip
	rm $(BINDIR)/vault.zip

depend: bin/mockgen bin/vault

go_generate: depend
	$(BINDIR)/mockgen -package kubernetes -source=pkg/kubernetes/kubernetes.go > pkg/kubernetes/kubernetes_mocks_test.go
	sed -i '1s/^/\/\/ Copyright Jetstack Ltd. See LICENSE for details.\n/' pkg/kubernetes/kubernetes_mocks_test.go

.builder_image:
	docker pull ${BUILD_IMAGE_NAME}


# Builder image targets
#######################
docker_%: .builder_image
	docker run -it \
		-v ${GOPATH}/src:/go/src \
		-v $(shell pwd):/go/src/${GO_PKG} \
		-w /go/src/${GO_PKG} \
		-e GOPATH=/go \
		${BUILD_IMAGE_NAME} \
		/bin/sh -c "make $*"

# Docker targets
################
docker_build:
	docker build -t $(REGISTRY)/$(IMAGE_NAME):$(BUILD_TAG) .

docker_push: docker_build
	set -e; \
		for tag in $(IMAGE_TAGS); do \
		docker tag $(REGISTRY)/$(IMAGE_NAME):$(BUILD_TAG) $(REGISTRY)/$(IMAGE_NAME):$${tag} ; \
		docker push $(REGISTRY)/$(IMAGE_NAME):$${tag}; \
		done
