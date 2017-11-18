load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "clientset.go",
        "doc.go",
    ],
    importpath = "kope.io/auth/pkg/client/clientset_generated/internalclientset",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/client/clientset_generated/internalclientset/typed/auth/internalversion:go_default_library",
        "//pkg/client/clientset_generated/internalclientset/typed/auth/v1alpha1:go_default_library",
        "//pkg/client/clientset_generated/internalclientset/typed/config/internalversion:go_default_library",
        "//pkg/client/clientset_generated/internalclientset/typed/config/v1alpha1:go_default_library",
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/k8s.io/client-go/discovery:go_default_library",
        "//vendor/k8s.io/client-go/rest:go_default_library",
        "//vendor/k8s.io/client-go/util/flowcontrol:go_default_library",
    ],
)
