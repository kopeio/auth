package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = glob(include = ["*.go"]),
    deps = [
        "@com_github_18F_hmacauth//:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "//pkg/apis/auth:go_default_library",
        "//pkg/cookie:go_default_library",
        "//pkg/cookie/proto:go_default_library",
        "//pkg/providers:go_default_library",
        "//pkg/tokenstore:go_default_library",
    ],
)
