package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_binary(
    name = "auth-api",
    library = ":go_default_library",
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    visibility = ["//visibility:private"],
    deps = [
        "//pkg/api/apiserver:go_default_library",
        "//pkg/client/clientset_generated/clientset:go_default_library",
        "//pkg/k8sauth:go_default_library",
        "//pkg/tokenstore:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_spf13_pflag//:go_default_library",
        "@io_k8s_apimachinery//pkg/util/wait:go_default_library",
        "@io_k8s_client_go//rest:go_default_library",
    ],
)
