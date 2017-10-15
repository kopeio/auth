package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["webhook.go"],
    importpath = "kope.io/auth/pkg/k8sauth",
    deps = [
        "//pkg/tokenstore:go_default_library",
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/k8s.io/api/authentication/v1beta1:go_default_library",
    ],
)
