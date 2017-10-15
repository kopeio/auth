load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "fake_authconfiguration.go",
        "fake_authprovider.go",
        "fake_componentconfig_client.go",
    ],
    importpath = "kope.io/auth/pkg/client/clientset_generated/internalclientset/typed/componentconfig/v1alpha1/fake",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/apis/componentconfig/v1alpha1:go_default_library",
        "//pkg/client/clientset_generated/internalclientset/typed/componentconfig/v1alpha1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/labels:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/schema:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/types:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/watch:go_default_library",
        "//vendor/k8s.io/client-go/rest:go_default_library",
        "//vendor/k8s.io/client-go/testing:go_default_library",
    ],
)
