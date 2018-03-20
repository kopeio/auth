package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "azure.go",
        "facebook.go",
        "github.go",
        "gitlab.go",
        "google.go",
        "internal_util.go",
        "linkedin.go",
        "myusa.go",
        "provider_data.go",
        "provider_default.go",
        "providers.go",
        "session_state.go",
    ],
    importpath = "kope.io/auth/pkg/providers",
)

go_test(
    name = "go_default_test",
    srcs = [
        "azure_test.go",
        "gitlab_test.go",
        "google_test.go",
        "internal_util_test.go",
        "linkedin_test.go",
        "myusa_test.go",
        "provider_default_test.go",
        "session_state_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//pkg/assert:go_default_library",
        "//pkg/cookie:go_default_library",
        "//pkg/cookie/pb:go_default_library",
    ],
)
