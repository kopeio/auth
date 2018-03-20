load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "pb_proto",
    srcs = ["cookie.proto"],
    visibility = ["//visibility:public"],
)

go_proto_library(
    name = "pb_go_proto",
    importpath = "kope.io/auth/pkg/oauth/pb",
    proto = ":pb_proto",
    visibility = ["//visibility:public"],
)

go_library(
    name = "go_default_library",
    embed = [":pb_go_proto"],
    importpath = "kope.io/auth/pkg/oauth/pb",
)
