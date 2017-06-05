package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "http.go",
        "kubeconfig.go",
        "login.go",
        "logout.go",
        "map_user.go",
        "oauth.go",
        "options.go",
        "portal.go",
        "rest.go",
    ],
    deps = [
        "//pkg/apis/auth/v1alpha1:go_default_library",
        "//pkg/configreader:go_default_library",
        "//pkg/keystore:go_default_library",
        "//pkg/oauth:go_default_library",
        "//pkg/oauth/session:go_default_library",
        "//pkg/tokenstore:go_default_library",
        "@com_github_ghodss_yaml//:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@io_k8s_apimachinery//pkg/types:go_default_library",
    ],
)
