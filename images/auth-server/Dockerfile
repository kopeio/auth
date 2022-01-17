FROM golang:1.17.4 AS builder

WORKDIR /workspace/auth
COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download
RUN CGO_ENABLED=0 go build k8s.io/client-go/rest sigs.k8s.io/controller-runtime/pkg/cache k8s.io/klog/v2 k8s.io/api/core/v1 k8s.io/client-go/kubernetes

COPY *.go .
COPY api/ api/
COPY pkg/ pkg/

RUN CGO_ENABLED=0 go build -v -o /workspace/auth-server .

FROM gcr.io/distroless/static-debian11:latest

COPY --from=builder /workspace/auth-server /bin/auth-server

ENTRYPOINT [ "/auth-server" ]