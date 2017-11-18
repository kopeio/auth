package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_docker//container:container.bzl",
    "container_image",
)

container_image(
    name = "auth-portal-assets",
    base = "@debian_base_amd64//image",
    data_path = "/webapp/portal/build/",
    directory = "/webapp",
    files = [
        "//webapp/portal:webapp",
    ],
)

container_image(
    name = "auth-portal",
    base = ":auth-portal-assets",
    cmd = ["/auth-portal"],
    files = [
        "//cmd/auth-portal",
    ],
)

container_image(
    name = "auth-api",
    base = "@debian_base_amd64//image",
    cmd = ["/user-api"],
    files = [
        "//cmd/auth-api",
    ],
)
