package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["cookies.go"],
    importpath = "kope.io/auth/pkg/cookie",
)

go_test(
    name = "go_default_test",
    srcs = ["cookies_test.go"],
    importpath = "kope.io/auth/pkg/cookie",
    library = ":go_default_library",
)
