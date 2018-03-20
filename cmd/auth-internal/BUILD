package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_binary(
    name = "auth-internal",
    embed = [":go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "kope.io/auth/cmd/auth-internal",
    visibility = ["//visibility:private"],
    deps = [
        "//pkg/api/apiserver:go_default_library",
        "//pkg/client/clientset_generated/clientset:go_default_library",
        "//pkg/k8sauth:go_default_library",
        "//pkg/tokenstore:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/wait:go_default_library",
        "//vendor/k8s.io/client-go/rest:go_default_library",
    ],
)
