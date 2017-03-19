package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "interface.go",
        "secrets.go",
        "tpr.go",
    ],
    deps = [
        "//pkg/apis/auth/v1alpha1:go_default_library",
        "//pkg/apis/auth/v1alpha1/client:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@io_k8s_apimachinery//pkg/apis/meta/v1:go_default_library",
        "@io_k8s_apimachinery//pkg/types:go_default_library",
        "@io_k8s_apimachinery//pkg/watch:go_default_library",
        "@io_k8s_client_go//pkg/apis/authentication/v1beta1:go_default_library",
        "@org_golang_x_crypto//bcrypt:go_default_library",
    ],
)
