DOCKER_TAG = latest
DOCKER_IMG = alisin69/msvc-entity-service:${DOCKER_TAG}

all: test

docker-push: docker-build
	docker push ${DOCKER_IMG}
.PHONY: docker-push

docker-build:
	docker build -t ${DOCKER_IMG} -f ./Dockerfile ..
.PHONY: docker-build

test: build
	go test -v -cover $(shell go list ./... | grep -v vendor)
.PHONY: test

build: codestyle
	go build
.PHONY: build

codestyle: gofmt golint govet
.PHONY: codestyle

gofmt:
	go fmt $(shell go list ./... | grep -v vendor)
.PHONY: gofmt

govet:
	go vet $(shell go list ./... | grep -v vendor)
.PHONY: govet

golint: $(GOPATH)/bin/golint
	golint $(shell go list ./... | grep -v vendor)
.PHONY: golint

dep: $(GOPATH)/bin/dep
	dep ensure -v
.PHONY: dep

$(GOPATH)/bin/golint:
	go get -u golang.org/x/lint/golint

$(GOPATH)/bin/dep:
	go get -u github.com/golang/dep
 