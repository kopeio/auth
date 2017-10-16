load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "default.go",
        "interface.go",
    ],
    importpath = "kope.io/auth/pkg/oauth/providers",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/apis/componentconfig/v1alpha1:go_default_library",
        "//pkg/oauth/session:go_default_library",
        "//vendor/github.com/golang/glog:go_default_library",
    ],
)
