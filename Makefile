# TODO: Move entirely to bazel?
.PHONY: images

DOCKER_REGISTRY?=kopeio
DOCKER_TAG=1.0.20170308

all: images

gofmt:
	gofmt -w -s cmd/
	gofmt -w -s pkg/

portal:
	cd webapp/portal; npm run build
	gzip --force --keep --best webapp/portal/public/bundle.js

push: images
	docker push ${DOCKER_REGISTRY}/auth-internal:${DOCKER_TAG}
	docker push ${DOCKER_REGISTRY}/auth-portal:${DOCKER_TAG}

images: portal
	bazel run //images:auth-internal ${DOCKER_REGISTRY}/auth-internal:${DOCKER_TAG}
	bazel run //images:auth-portal ${DOCKER_REGISTRY}/auth-portal:${DOCKER_TAG}
