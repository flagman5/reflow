load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["image_reference.go"],
    importmap = "grail/go/src/vendor/github.com/docker/docker/api/types/reference",
    importpath = "github.com/docker/docker/api/types/reference",
    visibility = ["//visibility:public"],
    deps = ["//go/src/vendor/github.com/docker/distribution/reference:go_default_library"],
)
