package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "htpasswd.go",
        "oauthproxy.go",
        "options.go",
        "templates.go",
        "version.go",
    ],
    importpath = "kope.io/auth/pkg/oauth2proxy",
    deps = [
        "//pkg/apis/auth/v1alpha1:go_default_library",
        "//pkg/cookie:go_default_library",
        "//pkg/cookie/pb:go_default_library",
        "//pkg/providers:go_default_library",
        "//pkg/tokenstore:go_default_library",
        "//vendor/github.com/18F/hmacauth:go_default_library",
        "//vendor/github.com/golang/glog:go_default_library",
    ],
)
