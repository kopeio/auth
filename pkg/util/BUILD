package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["api.go"],
    importpath = "kope.io/auth/pkg/util",
)

go_test(
    name = "go_default_test",
    srcs = ["api_test.go"],
    importpath = "kope.io/auth/pkg/util",
    library = ":go_default_library",
    deps = ["//pkg/assert:go_default_library"],
)
