http_archive(
    name = "io_bazel_rules_go",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.6.0/rules_go-0.6.0.tar.gz",
    sha256 = "ba6feabc94a5d205013e70792accb6cce989169476668fbaf98ea9b342e13b59",
)
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains", "go_repository")
go_rules_dependencies()
go_register_toolchains()


load("@io_bazel_rules_go//proto:def.bzl", "proto_register_toolchains")
proto_register_toolchains()

#=============================================================================

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

container_repositories()

container_pull(
    name = "debian_base_amd64",  # 'tag' is also supported, but digest is encouraged for reproducibility.
    digest = "sha256:f45c0513d335c6654ab72de6994f5899e65d87a24141e527b39ec5f503e22670",
    registry = "gcr.io",
    repository = "google_containers/debian-base-amd64",
    tag = "0.2",
)
