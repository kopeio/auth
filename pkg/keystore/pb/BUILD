package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

proto_library(
    name = "keystore_proto",
    srcs = ["keystore.proto"],
    visibility = ["//visibility:public"],
)

go_proto_library(
    name = "go_default_library",
    proto = ":keystore_proto",
    visibility = ["//visibility:public"],
)
