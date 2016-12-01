package(default_visibility = ["//visibility:public"])

load("@bazel_tools//tools/build_defs/docker:docker.bzl", "docker_build")

docker_build(
    name = "busybox",
    debs = [
        "@busybox_deb//file",
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
        "@libc_deb//file",
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
    name = "kauth-portal",
    base = ":busybox-ca-certificates",
    cmd = ["/kauth-portal"],
    files = [
        "//cmd/kauth-portal",
    ],
)

docker_build(
    name = "kauth-internal",
    base = ":busybox-libc",
    cmd = ["/kauth-internal"],
    files = [
        "//cmd/kauth-internal",
    ],
)
