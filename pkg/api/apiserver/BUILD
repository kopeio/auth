load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "server.go",
        "start.go",
        "utils.go",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/api/registry/authconfiguration:go_default_library",
        "//pkg/api/registry/authprovider:go_default_library",
        "//pkg/api/registry/user:go_default_library",
        "//pkg/apis/auth:go_default_library",
        "//pkg/apis/auth/install:go_default_library",
        "//pkg/apis/auth/v1alpha1:go_default_library",
        "//pkg/apis/componentconfig:go_default_library",
        "//pkg/apis/componentconfig/install:go_default_library",
        "//pkg/apis/componentconfig/v1alpha1:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_spf13_pflag//:go_default_library",
        "@io_k8s_apimachinery//pkg/apimachinery/announced:go_default_library",
        "@io_k8s_apimachinery//pkg/apimachinery/registered:go_default_library",
        "@io_k8s_apimachinery//pkg/apis/meta/v1:go_default_library",
        "@io_k8s_apimachinery//pkg/runtime:go_default_library",
        "@io_k8s_apimachinery//pkg/runtime/schema:go_default_library",
        "@io_k8s_apimachinery//pkg/runtime/serializer:go_default_library",
        "@io_k8s_apimachinery//pkg/version:go_default_library",
        "@io_k8s_apiserver//pkg/registry/rest:go_default_library",
        "@io_k8s_apiserver//pkg/server:go_default_library",
        "@io_k8s_apiserver//pkg/server/options:go_default_library",
        "@io_k8s_apiserver//pkg/storage/storagebackend:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "register_test.go",
        "scheme_test.go",
    ],
    library = ":go_default_library",
    deps = [
        ":go_default_library",
        "//pkg/apis/componentconfig/v1alpha1:go_default_library",
        "@io_k8s_apimachinery//pkg/api/testing:go_default_library",
        "@io_k8s_apimachinery//pkg/runtime:go_default_library",
    ],
)
