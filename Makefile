DOCKER_REGISTRY?=$(shell whoami)
DOCKER_TAG?=latest

all: images

gofmt:
	gofmt -w -s cmd/
	gofmt -w -s pkg/

portal:
	cd webapp/portal; npm run build
	gzip --force --keep --best webapp/portal/build/static/js/main.*.js

portal-image:
	bazel run //images:auth-portal
	docker tag bazel/images:auth-portal ${DOCKER_REGISTRY}/auth-portal:${DOCKER_TAG}

portal-push: portal-image
	docker push ${DOCKER_REGISTRY}/auth-portal:${DOCKER_TAG}

portal-bounce:
	kubectl delete pod -n kopeio-auth -l app=auth-portal

api-image:
	bazel run //images:auth-api
	docker tag bazel/images:auth-api ${DOCKER_REGISTRY}/auth-api:${DOCKER_TAG}

api-push: api-image
	docker push ${DOCKER_REGISTRY}/auth-api:${DOCKER_TAG}

api-bounce:
	kubectl delete pod -n kopeio-auth -l app=auth-api

use-dev-images:
	kubectl set image ds -n kopeio-auth auth-api auth-api=${DOCKER_REGISTRY}/auth-api:${DOCKER_TAG}
	kubectl set image deployment -n kopeio-auth auth-portal auth-portal=${DOCKER_REGISTRY}/auth-portal:${DOCKER_TAG}

push: portal-push api-push
	echo "pushed images"

.PHONY: images
images: portal-image api-image
	echo "built images"

dep:
	dep ensure
	find vendor -name "BUILD" -delete
	find vendor -name "BUILD.bazel" -delete
	bazel run //:gazelle -- -proto disable

.PHONY: gazelle
gazelle:
	bazel run //:gazelle
	git checkout -- vendor/
	git clean -df vendor/

goimports:
	goimports -w cmd/ pkg/

# -----------------------------------------------------
# api machinery regenerate

apimachinery:
	#./hack/make-apimachinery.sh
	${GOPATH}/bin/conversion-gen --skip-unsafe=true --input-dirs kope.io/auth/pkg/apis/auth/v1alpha1 --v=0  --output-file-base=zz_generated.conversion --go-header-file hack/boilerplate/boilerplate.go.txt --extra-peer-dirs=k8s.io/api/core/v1,k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/conversion,k8s.io/apimachinery/pkg/runtime
	${GOPATH}/bin/conversion-gen --skip-unsafe=true --input-dirs kope.io/auth/pkg/apis/componentconfig/v1alpha1 --v=0  --output-file-base=zz_generated.conversion --go-header-file hack/boilerplate/boilerplate.go.txt --extra-peer-dirs=k8s.io/api/core/v1,k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/conversion,k8s.io/apimachinery/pkg/runtime
	${GOPATH}/bin/defaulter-gen --input-dirs kope.io/auth/pkg/apis/auth/v1alpha1 --v=0  --output-file-base=zz_generated.defaults --go-header-file hack/boilerplate/boilerplate.go.txt
	${GOPATH}/bin/defaulter-gen --input-dirs kope.io/auth/pkg/apis/componentconfig/v1alpha1 --v=0  --output-file-base=zz_generated.defaults --go-header-file hack/boilerplate/boilerplate.go.txt
	${GOPATH}/bin/deepcopy-gen --input-dirs kope.io/auth/pkg/apis/auth --v=0  --output-file-base=zz_generated.deepcopy --go-header-file hack/boilerplate/boilerplate.go.txt
	${GOPATH}/bin/deepcopy-gen --input-dirs kope.io/auth/pkg/apis/auth/v1alpha1 --v=0  --output-file-base=zz_generated.deepcopy --go-header-file hack/boilerplate/boilerplate.go.txt
	${GOPATH}/bin/deepcopy-gen --input-dirs kope.io/auth/pkg/apis/componentconfig --v=0  --output-file-base=zz_generated.deepcopy --go-header-file hack/boilerplate/boilerplate.go.txt
	${GOPATH}/bin/deepcopy-gen --input-dirs kope.io/auth/pkg/apis/componentconfig/v1alpha1 --v=0  --output-file-base=zz_generated.deepcopy --go-header-file hack/boilerplate/boilerplate.go.txt
	${GOPATH}/bin/client-gen  --input-base kope.io/auth/pkg/apis --input="auth/,auth/v1alpha1,componentconfig/,componentconfig/v1alpha1" --clientset-path kope.io/auth/pkg/client/clientset_generated/ --go-header-file hack/boilerplate/boilerplate.go.txt
	${GOPATH}/bin/client-gen  --clientset-name="clientset" --input-base kope.io/auth/pkg/apis --input="auth/v1alpha1,componentconfig/v1alpha1" --clientset-path kope.io/auth/pkg/client/clientset_generated/ --go-header-file hack/boilerplate/boilerplate.go.txt
