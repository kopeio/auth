package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = glob(include = ["*.go"]),
    deps = [
        "@com_github_ghodss_yaml//:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@io_k8s_client_go//kubernetes:go_default_library",
        "@io_k8s_client_go//pkg/apis/authentication/v1beta1:go_default_library",
        "@io_k8s_client_go//rest:go_default_library",
        "//pkg/apis/auth:go_default_library",
        "//pkg/oauth2proxy:go_default_library",
        "//pkg/tokenstore:go_default_library",
    ],
)
