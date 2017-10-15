http_archive(
    name = "io_bazel_rules_go",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.6.0/rules_go-0.6.0.tar.gz",
    sha256 = "ba6feabc94a5d205013e70792accb6cce989169476668fbaf98ea9b342e13b59",
)
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains", "go_repository")
go_rules_dependencies()
go_register_toolchains()

#============================================================================
git_repository(
    name = "io_bazel_rules_docker",
    remote = "https://github.com/bazelbuild/rules_docker.git",
    tag = "v0.3.0",
)

load(
    "@io_bazel_rules_docker//container:container.bzl",
    "container_pull",
    container_repositories = "repositories",
)

# This is NOT needed when going through the language lang_image
# "repositories" function(s).
container_repositories()

#=============================================================================

load("@io_bazel_rules_go//proto:def.bzl", "proto_register_toolchains")
proto_register_toolchains()

#=============================================================================

# for building docker base images
debs = (
    (
        "deb_busybox",
        "5f81f140777454e71b9e5bfdce9c89993de5ddf4a7295ea1cfda364f8f630947",
        "http://ftp.us.debian.org/debian/pool/main/b/busybox/busybox-static_1.22.0-19+b3_amd64.deb",
    ),
    (
        "deb_libc",
        "b3f7278d80d5d0dc428fe92309bbc0e0a1ed665548a9f660663c1e1151335ae9",
        "http://ftp.us.debian.org/debian/pool/main/g/glibc/libc6_2.24-11+deb9u1_amd64.deb",
    ),
)

[http_file(
    name = name,
    sha256 = sha256,
    url = url,
) for name, sha256, url in debs]
