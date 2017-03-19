load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["api.go"],
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_golang_glog//:go_default_library",
        "@io_k8s_apimachinery//pkg/apimachinery/announced:go_default_library",
        "@io_k8s_apimachinery//pkg/apimachinery/registered:go_default_library",
        "@io_k8s_apimachinery//pkg/runtime:go_default_library",
        "@io_k8s_apimachinery//pkg/runtime/serializer:go_default_library",
    ],
)
