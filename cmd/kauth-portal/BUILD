package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_prefix", "go_binary")

go_binary(
    name = "kauth-portal",
    srcs = glob(
        include = ["*.go"],
    ),
    deps = [
        "@com_github_golang_glog//:go_default_library",
        "@io_k8s_client_go//kubernetes:go_default_library",
        "@io_k8s_client_go//pkg/util/yaml:go_default_library",
        "@io_k8s_client_go//rest:go_default_library",
        "//pkg/portal:go_default_library",
        "//pkg/tokenstore:go_default_library",
    ],
)
