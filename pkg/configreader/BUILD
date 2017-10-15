load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["reader.go"],
    importpath = "kope.io/auth/pkg/configreader",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/apis/componentconfig/v1alpha1:go_default_library",
        "//pkg/client/clientset_generated/clientset:go_default_library",
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/fields:go_default_library",
        "//vendor/k8s.io/client-go/tools/cache:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["reader_test.go"],
    importpath = "kope.io/auth/pkg/configreader",
    library = ":go_default_library",
    deps = [
        "//pkg/api/apiserver:go_default_library",
        "//pkg/apis/componentconfig:go_default_library",
        "//pkg/apis/componentconfig/v1alpha1:go_default_library",
        "//vendor/k8s.io/api/core/v1:go_default_library",
    ],
)
