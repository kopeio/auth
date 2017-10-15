package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

proto_library(
    name = "cookie_proto",
    srcs = ["cookie.proto"],
    visibility = ["//visibility:public"],
)

go_proto_library(
    name = "go_default_library",
    proto = ":cookie_proto",
    visibility = ["//visibility:public"],
)
