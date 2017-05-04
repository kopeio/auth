load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "clientset_generated.go",
        "doc.go",
        "register.go",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/apis/auth/install:go_default_library",
        "//pkg/client/clientset_generated/internalclientset:go_default_library",
        "//pkg/client/clientset_generated/internalclientset/typed/auth/internalversion:go_default_library",
        "//pkg/client/clientset_generated/internalclientset/typed/auth/internalversion/fake:go_default_library",
        "//pkg/client/clientset_generated/internalclientset/typed/auth/v1alpha1:go_default_library",
        "//pkg/client/clientset_generated/internalclientset/typed/auth/v1alpha1/fake:go_default_library",
        "@io_k8s_apimachinery//pkg/apimachinery/announced:go_default_library",
        "@io_k8s_apimachinery//pkg/apimachinery/registered:go_default_library",
        "@io_k8s_apimachinery//pkg/apis/meta/v1:go_default_library",
        "@io_k8s_apimachinery//pkg/runtime:go_default_library",
        "@io_k8s_apimachinery//pkg/runtime/schema:go_default_library",
        "@io_k8s_apimachinery//pkg/runtime/serializer:go_default_library",
        "@io_k8s_apimachinery//pkg/watch:go_default_library",
        "@io_k8s_client_go//discovery:go_default_library",
        "@io_k8s_client_go//discovery/fake:go_default_library",
        "@io_k8s_client_go//testing:go_default_library",
    ],
)
