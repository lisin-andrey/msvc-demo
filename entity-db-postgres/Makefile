DOCKER_TAG = latest
DOCKER_IMG = alisin69/msvc-pgdb:${DOCKER_TAG}

all: docker-build

docker-push: docker-build
	docker push ${DOCKER_IMG}
.PHONY: docker-push

docker-build:
	docker build -t ${DOCKER_IMG} .
.PHONY: docker-build
