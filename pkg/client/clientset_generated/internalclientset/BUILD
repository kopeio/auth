load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "clientset.go",
        "doc.go",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/client/clientset_generated/internalclientset/typed/auth/internalversion:go_default_library",
        "//pkg/client/clientset_generated/internalclientset/typed/auth/v1alpha1:go_default_library",
        "//pkg/client/clientset_generated/internalclientset/typed/componentconfig/internalversion:go_default_library",
        "//pkg/client/clientset_generated/internalclientset/typed/componentconfig/v1alpha1:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@io_k8s_client_go//discovery:go_default_library",
        "@io_k8s_client_go//rest:go_default_library",
        "@io_k8s_client_go//util/flowcontrol:go_default_library",
    ],
)
