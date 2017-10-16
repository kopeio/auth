load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["session.go"],
    importpath = "kope.io/auth/pkg/oauth/session",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/oauth/pb:go_default_library",
        "//vendor/github.com/golang/protobuf/proto:go_default_library",
    ],
)
