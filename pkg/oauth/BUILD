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
    deps = [
        "//pkg/configreader:go_default_library",
        "//pkg/keystore:go_default_library",
        "//pkg/oauth/pb:go_default_library",
        "//pkg/oauth/providers:go_default_library",
        "//pkg/oauth/providers/google:go_default_library",
        "//pkg/oauth/session:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@io_k8s_apimachinery//pkg/types:go_default_library",
    ],
)
