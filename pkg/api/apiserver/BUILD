load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "server.go",
        "start.go",
        "utils.go",
    ],
    importpath = "kope.io/auth/pkg/api/apiserver",
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
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/github.com/spf13/pflag:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apimachinery/announced:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apimachinery/registered:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/schema:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/serializer:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/version:go_default_library",
        "//vendor/k8s.io/apiserver/pkg/registry/rest:go_default_library",
        "//vendor/k8s.io/apiserver/pkg/server:go_default_library",
        "//vendor/k8s.io/apiserver/pkg/server/options:go_default_library",
        "//vendor/k8s.io/apiserver/pkg/storage/storagebackend:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "register_test.go",
        "scheme_test.go",
    ],
    importpath = "kope.io/auth/pkg/api/apiserver",
    library = ":go_default_library",
    deps = [
        ":go_default_library",
        "//pkg/apis/componentconfig/v1alpha1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/api/testing:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
    ],
)
