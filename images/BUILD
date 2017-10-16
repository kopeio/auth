package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_docker//docker:docker.bzl",
    "docker_build",
)

docker_build(
    name = "busybox",
    debs = [
        "@deb_busybox//file",
    ],
    symlinks = {
        "/bin/sh": "/bin/busybox",
        "/usr/bin/busybox": "/bin/busybox",
        "/usr/sbin/busybox": "/bin/busybox",
        "/sbin/busybox": "/bin/busybox",
    },
)

docker_build(
    name = "busybox-libc",
    base = ":busybox",
    debs = [
        "@deb_libc//file",
    ],
)

docker_build(
    name = "busybox-ca-certificates",
    base = ":busybox-libc",
    cmd = ["/oauth2-proxy"],
    data_path = "/images",
    directory = "/etc/ssl/certs",
    files = [
        "//images:ca-certificates.crt",
    ],
)

docker_build(
    name = "auth-portal-assets",
    base = ":busybox-ca-certificates",
    data_path = "/webapp/portal/build/",
    directory = "/webapp",
    files = [
        "//webapp/portal:webapp",
    ],
)

docker_build(
    name = "auth-portal",
    base = ":auth-portal-assets",
    cmd = ["/auth-portal"],
    files = [
        "//cmd/auth-portal",
    ],
)

docker_build(
    name = "auth-api",
    base = ":busybox-libc",
    cmd = ["/user-api"],
    files = [
        "//cmd/auth-api",
    ],
)
