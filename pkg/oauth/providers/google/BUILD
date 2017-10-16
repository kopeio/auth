load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["google.go"],
    importpath = "kope.io/auth/pkg/oauth/providers/google",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/apis/componentconfig/v1alpha1:go_default_library",
        "//pkg/oauth/pb:go_default_library",
        "//pkg/oauth/providers:go_default_library",
        "//pkg/oauth/session:go_default_library",
        "//vendor/github.com/golang/glog:go_default_library",
    ],
)
