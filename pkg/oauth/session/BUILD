load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["session.go"],
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/oauth/pb:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
    ],
)
