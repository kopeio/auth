# TODO: Move entirely to bazel?
.PHONY: images

DOCKER_REGISTRY?=kopeio
DOCKER_TAG=1.0.20170204

all: images

gofmt:
	gofmt -w -s cmd/
	gofmt -w -s pkg/

push: images
	docker push ${DOCKER_REGISTRY}/auth-internal:${DOCKER_TAG}
	docker push ${DOCKER_REGISTRY}/auth-portal:${DOCKER_TAG}

images:
	bazel run //images:auth-internal ${DOCKER_REGISTRY}/auth-internal:${DOCKER_TAG}
	bazel run //images:auth-portal ${DOCKER_REGISTRY}/auth-portal:${DOCKER_TAG}
