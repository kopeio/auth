load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["google.go"],
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/apis/componentconfig/v1alpha1:go_default_library",
        "//pkg/oauth/pb:go_default_library",
        "//pkg/oauth/providers:go_default_library",
        "//pkg/oauth/session:go_default_library",
        "@com_github_golang_glog//:go_default_library",
    ],
)
