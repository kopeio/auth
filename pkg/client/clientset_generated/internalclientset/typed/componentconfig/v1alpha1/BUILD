load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "authconfiguration.go",
        "componentconfig_client.go",
        "doc.go",
        "generated_expansion.go",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/apis/componentconfig/v1alpha1:go_default_library",
        "//pkg/client/clientset_generated/internalclientset/scheme:go_default_library",
        "@io_k8s_apimachinery//pkg/apis/meta/v1:go_default_library",
        "@io_k8s_apimachinery//pkg/runtime/serializer:go_default_library",
        "@io_k8s_apimachinery//pkg/types:go_default_library",
        "@io_k8s_apimachinery//pkg/watch:go_default_library",
        "@io_k8s_client_go//rest:go_default_library",
    ],
)
