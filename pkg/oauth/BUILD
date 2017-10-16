package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "authenticate.go",
        "cookies.go",
        "http.go",
        "providers.go",
        "state.go",
    ],
    importpath = "kope.io/auth/pkg/oauth",
    deps = [
        "//pkg/configreader:go_default_library",
        "//pkg/keystore:go_default_library",
        "//pkg/oauth/pb:go_default_library",
        "//pkg/oauth/providers:go_default_library",
        "//pkg/oauth/providers/google:go_default_library",
        "//pkg/oauth/session:go_default_library",
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/github.com/golang/protobuf/proto:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/types:go_default_library",
    ],
)
