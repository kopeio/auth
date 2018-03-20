package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_binary(
    name = "auth-portal",
    embed = [":go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "kope.io/auth/cmd/auth-portal",
    visibility = ["//visibility:private"],
    deps = [
        "//pkg/client/clientset_generated/clientset:go_default_library",
        "//pkg/configreader:go_default_library",
        "//pkg/keystore:go_default_library",
        "//pkg/portal:go_default_library",
        "//pkg/tokenstore:go_default_library",
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes:go_default_library",
        "//vendor/k8s.io/client-go/rest:go_default_library",
    ],
)
