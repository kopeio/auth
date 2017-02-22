package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = glob(include = ["*.go"]),
    deps = [
        "@com_github_golang_glog//:go_default_library",
        #        "@com_github_bitly_go-simplejson//:go_default_library",
        #        "@com_github_bmizerany_assert//:go_default_library",
        "@org_golang_google_api//admin/directory/v1:go_default_library",
        "@org_golang_x_oauth2//:go_default_library",
        "@org_golang_x_oauth2//google:go_default_library",
        "//pkg/assert:go_default_library",
        "//pkg/util:go_default_library",
        "//pkg/cookie:go_default_library",
        "//pkg/cookie/proto:go_default_library",
    ],
)
