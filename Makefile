REGISTRY := justinsb
TAG := latest

push:
	docker buildx build --push -t ${REGISTRY}/auth-server:${TAG} -f images/auth-server/Dockerfile .

bounce:
	kubectl delete pod -n kopeio-auth-system --all

logs:
	kubectl logs -n kopeio-auth-system -l app=auth-server

generate:
	go generate ./...

.PHONY: protoc
protoc:
	protoc --go_out=. --go_opt=paths=source_relative ./pkg/keystore/pb/keystore.proto ./pkg/session/pb/cookie.proto ./pkg/oauth/pb/state.proto